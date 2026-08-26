package controllers

import (
	"forum/src/models"
	"net/http"
	"time"
)

func SetCookie(user models.UserType) *http.Cookie {
	return &http.Cookie{
		Name:     "__Host-FRMSessionID",
		Value:    user.SessionId,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSite(http.SameSiteStrictMode),
	}
}

func UnsetCookie() *http.Cookie {
	return &http.Cookie{
		Name:     "__Host-FRMSessionID",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSite(http.SameSiteStrictMode),
	}
}
