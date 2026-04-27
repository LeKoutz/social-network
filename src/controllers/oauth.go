package controllers

import (
	"encoding/json"
	"errors"
	"forum/src/models"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

type oauthConfig struct {
	ClientID, ClientSecret, RedirectURL, AuthURL, TokenURL string
	Scopes []string
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


var (
	googleOAuthConf = &oauthConfig{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		Scopes:       []string{"email", "profile"},
		AuthURL:      "https://accounts.google.com/o/oauth2/v2/auth",
		TokenURL:     "https://oauth2.googleapis.com/token",
	}
	githubOAuthConf = &oauthConfig{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/auth/github/callback",
		Scopes:       []string{"user:email", "read:user"},
		AuthURL:      "https://github.com/login/oauth/authorize",
		TokenURL:	  "https://github.com/login/oauth/access_token",
	}
)

func handleGoogleLogin(data models.ResponseStruct) {
	state, _ := uuid.NewV4()
	url := googleOAuthConf.AuthCodeURL(state.String()) + "&prompt=consent"
	http.Redirect(data.Response, data.Request, url, http.StatusTemporaryRedirect)
}

func handleGitHubLogin(data models.ResponseStruct) {
	state, _ := uuid.NewV4()
	url := githubOAuthConf.AuthCodeURL(state.String()) + "&prompt=consent"
	http.Redirect(data.Response, data.Request, url, http.StatusTemporaryRedirect)
}

func handleGoogleCallback(data models.ResponseStruct) {
	token, err := googleOAuthConf.Exchange(data.Request.Context(), data.Request.URL.Query().Get("code"))
	if err != nil {
		data.Error.Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	resp, err := googleOAuthConf.Client(data.Request.Context(), token).Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil || resp == nil {
		data.Error.Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	defer resp.Body.Close()

	var info struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	json.NewDecoder(resp.Body).Decode(&info)

	username := strings.ReplaceAll(info.Name, " ", "_")
	createOrLoginUser(data, "google", info.Email, username)
}

func handleGitHubCallback(data models.ResponseStruct) {
	token, err := githubOAuthConf.Exchange(data.Request.Context(), data.Request.URL.Query().Get("code"))
	if err != nil {
		data.Error.Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	resp, err := githubOAuthConf.Client(data.Request.Context(), token).Get("https://api.github.com/user")
	if err != nil || resp == nil {
		data.Error.Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	defer resp.Body.Close()

	var info struct {
		ID    string `json:"id"`
		Login string `json:"login"`
	}
	json.NewDecoder(resp.Body).Decode(&info)

	emailResp, _ := githubOAuthConf.Client(data.Request.Context(), token).Get("https://api.github.com/user/emails")
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
	createOrLoginUser(data, "github", email, username)
}

func createOrLoginUser(data models.ResponseStruct, provider, email, username string) {
	var user models.User
	var err error
	if !models.IsEmailRegistered(email) {
		user.Username = username
		user.Email = email
		user.OAuthProvider = provider
		err := user.AddOAuth()
		if err != nil {
			data.Error.Consume(err).LogAndRespondError(data.Response, data.User)
			return
		}
	}
	sessionValue, err := uuid.NewV4()
	if err != nil {
		data.User = models.GetGuestUser()
		data.Error.Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	data.User, err = models.GetUserByOAuthProviderAndEmail(provider, email)
	if err != nil {
		data.User = models.GetGuestUser()
		if errors.Is(err, models.ErrorNoRows){
			data.Error.Consume(models.ErrorEmailNotFoundForOAuth).LogAndRespondError(data.Response, data.User)
			return
		}
		data.Error.Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	data.User.LoggedIn = true
	err = data.User.SetUserSession(sessionValue.String())
	if err != nil {
		data.User = models.GetGuestUser()
		data.Error.Consume(err).LogAndRespondError(data.Response, data.User)
		return
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
	http.SetCookie(data.Response, cookie)

	// http.Redirect(data.Response, data.Request, "/", http.StatusSeeOther)
	//
	// Hack: cause if we redirect, it somehow doesn't read the cookie after the
	// redirect. The cookie is set, though. It just doesn't leave the browser at
	// this point. So, instead, we return them to the Index() controller without
	// redirection... weird... the exact same flow goes fine on attemptLogin()
	// and it *does* redirect. Maybe because we redirected from somewhere else?
	Index(data)
}
