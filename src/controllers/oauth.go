package controllers

import (
	"encoding/json"
	"forum/src/models"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

var (
	googleOAuthConf = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}
	githubOAuthConf = &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/auth/github/callback",
		Scopes:       []string{"user:email", "read:user"},
		Endpoint:     github.Endpoint,
	}
)

func handleGoogleLogin(data models.ResponseStruct) {
	state, _ := uuid.NewV4()
	http.Redirect(data.Response, data.Request, googleOAuthConf.AuthCodeURL(state.String()), http.StatusTemporaryRedirect)
}

func handleGitHubLogin(data models.ResponseStruct) {
	state, _ := uuid.NewV4()
	http.Redirect(data.Response, data.Request, githubOAuthConf.AuthCodeURL(state.String()), http.StatusTemporaryRedirect)
}

func handleGoogleCallback(data models.ResponseStruct) {
	token, err := googleOAuthConf.Exchange(data.Request.Context(), data.Request.URL.Query().Get("code"))
	if err != nil {
		data.Error.Consume(err).LogError()
		return
	}
	resp, err := googleOAuthConf.Client(data.Request.Context(), token).Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil || resp == nil {
		data.Error.Consume(err).LogError()
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
	createOrLoginUser(data, "google", info.ID, info.Email, username)
}

func handleGitHubCallback(data models.ResponseStruct) {
	token, err := githubOAuthConf.Exchange(data.Request.Context(), data.Request.URL.Query().Get("code"))
	if err != nil {
		data.Error.Consume(err).LogError()
		return
	}
	resp, err := githubOAuthConf.Client(data.Request.Context(), token).Get("https://api.github.com/user")
	if err != nil || resp == nil {
		data.Error.Consume(err).LogError()
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
	createOrLoginUser(data, "github", info.ID, email, username)
}

func createOrLoginUser(data models.ResponseStruct, provider, oauthID, email, username string) {
	user, err := models.GetUserByOAuth(provider, oauthID)
	if err != nil {
		data.Error.Consume(err).LogError()
		return
	}
	if user.Id == 0 {
		user = models.User{
			Username:      username,
			Email:         email,
			OAuthProvider: provider,
			OAuthId:       oauthID,
		}
		err := user.AddOAuth()
		if err != nil {
			data.Error.Consume(err).LogError()
			return
		}
		user, err = models.GetUserByOAuth(provider, oauthID)
		if err != nil {
			data.Error.Consume(err).LogError()
			return
		}
	}
	session, _ := uuid.NewV4()
	user.SetUserSession(session.String())
	cookie := &http.Cookie{
		Name:     "__Host-FRMSessionID",
		Value:    session.String(),
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSite(http.SameSiteStrictMode),
	}
	http.SetCookie(data.Response, cookie)
	http.Redirect(data.Response, data.Request, "/", http.StatusSeeOther)
}
