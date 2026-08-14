package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"forum/src/ferror"
	"forum/src/models"
	"forum/src/state"
	"forum/src/utils"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

type oauthConfig struct {
	ClientID, ClientSecret, RedirectURL, AuthURL, TokenURL string
	Scopes                                                 []string
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
}

type BearerTransport struct {
	Token string
}

func (c *oauthConfig) AuthCodeURL(state string) string {
	params := url.Values{
		"client_id":     {c.ClientID},
		"redirect_uri":  {c.RedirectURL},
		"response_type": {"code"},
		"scope":         {strings.Join(c.Scopes, " ")},
		"state":         {state},
	}
	return c.AuthURL + "?" + params.Encode()
}

func (c *oauthConfig) Exchange(ctx context.Context, code string) (string, error) {
	data := url.Values{
		"client_id":     {c.ClientID},
		"client_secret": {c.ClientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {c.RedirectURL},
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tr tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return "", err
	}
	if tr.AccessToken == "" {
		return "", ferror.ErrorAccessToken
	}

	return tr.AccessToken, nil
}

func (b *BearerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+b.Token)
	req.Header.Set("Accept", "application/json")
	return http.DefaultTransport.RoundTrip(req)
}

func (c *oauthConfig) Client(ctx context.Context, token string) *http.Client {
	return &http.Client{
		Transport: &BearerTransport{Token: token},
	}
}

var (
	googleOAuthConf = &oauthConfig{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("BASE_URL") + "auth/google/callback",
		Scopes:       []string{"email", "profile"},
		AuthURL:      "https://accounts.google.com/o/oauth2/v2/auth",
		TokenURL:     "https://oauth2.googleapis.com/token",
	}
	githubOAuthConf = &oauthConfig{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("BASE_URL") + "auth/github/callback",
		Scopes:       []string{"user:email", "read:user"},
		AuthURL:      "https://github.com/login/oauth/authorize",
		TokenURL:     "https://github.com/login/oauth/access_token",
	}
)

func HandleOAuthLogin(data state.StateController, provider string) {
	state, _ := uuid.NewV4()

	cookie := &http.Cookie{
		Name:     "__Host-FRMState",
		Value:    state.String(),
		Path:     "/",
		Expires:  time.Now().Add(10 * time.Minute),
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(*data.EditResponse(), cookie)

	var url string
	switch provider {
	case "google":
		url = googleOAuthConf.AuthCodeURL(state.String()) + "&prompt=consent"
	case "github":
		url = githubOAuthConf.AuthCodeURL(state.String()) + "&prompt=consent"
	}
	http.Redirect(*data.EditResponse(), data.GetRequest(), url, http.StatusTemporaryRedirect)
}

func OAuthGoogleCallback(data state.StateController) error {
	token, err := googleOAuthConf.Exchange(data.GetRequest().Context(), data.GetRequest().URL.Query().Get("code"))
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	resp, err := googleOAuthConf.Client(data.GetRequest().Context(), token).Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil || resp == nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	defer resp.Body.Close()

	var info struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	json.NewDecoder(resp.Body).Decode(&info)

	username := strings.ReplaceAll(info.Name, " ", "_")
	return createOrLoginUser(data, "google", info.Email, username)
}

func OAuthGitHubCallback(data state.StateController) error {
	token, err := githubOAuthConf.Exchange(data.GetRequest().Context(), data.GetRequest().URL.Query().Get("code"))
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	resp, err := githubOAuthConf.Client(data.GetRequest().Context(), token).Get("https://api.github.com/user")
	if err != nil || resp == nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	defer resp.Body.Close()

	var info struct {
		ID    string `json:"id"`
		Login string `json:"login"`
	}
	json.NewDecoder(resp.Body).Decode(&info)

	emailResp, _ := githubOAuthConf.Client(data.GetRequest().Context(), token).Get("https://api.github.com/user/emails")
	var email string
	if emailResp != nil {
		defer emailResp.Body.Close()
		var emails []struct {
			Email   string `json:"email"`
			Primary bool   `json:"primary"`
		}
		json.NewDecoder(emailResp.Body).Decode(&emails)
		for _, e := range emails {
			if e.Primary {
				email = e.Email
				break
			}
		}
	}

	username := strings.ReplaceAll(info.Login, " ", "_")
	return createOrLoginUser(data, "github", email, username)
}

func createOrLoginUser(data state.StateController, provider, email, username string) error {
	var err error
	if !models.IsEmailRegistered(email) {
		data.EditUser().Username = username
		data.EditUser().Email = email
		data.EditUser().OAuthProvider = provider
		err := data.EditUser().AddOAuth()
		if err != nil {
			return errors.Join(utils.GetFunctionName(), err)
		}
	} else {
		data.EditUser().Email = email
		data.EditUser().OAuthProvider = provider
	}
	sessionValue, err := uuid.NewV4()
	if err != nil {
		data.SetUser(models.GetGuestUser())
		return errors.Join(utils.GetFunctionName(), err)
	}
	err = data.EditUser().GetUserByOAuthProviderAndEmail()
	if err != nil {
		data.SetUser(models.GetGuestUser())
		if errors.Is(err, ferror.ErrorNoRows) {
			return ferror.ErrorEmailNotFoundForOAuth
		}
		return errors.Join(utils.GetFunctionName(), err)
	}
	data.EditUser().LoggedIn = true
	err = data.EditUser().SetUserSession(sessionValue.String())
	if err != nil {
		data.SetUser(models.GetGuestUser())
		return errors.Join(utils.GetFunctionName(), err)
	}
	err = data.EditUser().GetNotifications()
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	cookie := &http.Cookie{
		Name:     "__Host-FRMSessionID",
		Value:    sessionValue.String(),
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSite(http.SameSiteStrictMode),
	}
	http.SetCookie(*data.EditResponse(), cookie)

	// Because if we redirect, it somehow doesn't read the cookie after the
	// redirect. The cookie is set, though. It just doesn't leave the browser at
	// this point. So, instead, we return them to the Index() controller without
	// redirection... The difference with internal auth attemptLogin() where we
	// do redirection successfully is because we are doing it from the same
	// origin. Here, we land from external page.
	Index(data)
	return nil
}
