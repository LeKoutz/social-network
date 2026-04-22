package controllers

import (
	"forum/src/models"
	"forum/src/utils"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func Routes(data models.ResponseStruct) {
	uri, err := url.ParseRequestURI(data.Request.RequestURI)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	utils.LogDebug(uri)
	switch {
	case strings.HasPrefix(data.Request.RequestURI, "/auth/google/callback"):
		handleGoogleCallback(data)
	case strings.HasPrefix(data.Request.RequestURI, "/auth/google"):
		handleOAuthLogin(data, "google")
	case strings.HasPrefix(data.Request.RequestURI, "/auth/github/callback"):
		handleGitHubCallback(data)
	case strings.HasPrefix(data.Request.RequestURI, "/auth/github"):
		handleOAuthLogin(data, "github")
	case strings.HasPrefix(data.Request.RequestURI, "/category/view/"):
		if data.Request.Method == http.MethodGet {
			showCategory(data)
		} else {
			(&models.Error{}).Consume(models.ErrorMethodNotAllowed).LogAndRespondError(data.Response, data.User)
		}
	case strings.HasPrefix(data.Request.RequestURI, "/comment/react"):
		handleCommentReaction(data)
	case strings.HasPrefix(data.Request.RequestURI, "/comment/create"):
		handleCommentCreate(data)
	case strings.Compare(data.Request.RequestURI, "/categories") == 0:
		if data.Request.Method == http.MethodGet {
			showCategories(data)
		} else {
			(&models.Error{}).Consume(models.ErrorMethodNotAllowed).LogAndRespondError(data.Response, data.User)
		}
	case strings.Compare(data.Request.RequestURI, "/posts") == 0:
		if data.Request.Method == http.MethodGet {
			showPosts(data)
		} else {
			(&models.Error{}).Consume(models.ErrorMethodNotAllowed).LogAndRespondError(data.Response, data.User)
		}
	case strings.Compare(data.Request.RequestURI, "/post/create") == 0:
		handlePostCreate(data)
	case strings.Compare(data.Request.RequestURI, "/post/react") == 0:
		handlePostReaction(data)
	case strings.HasPrefix(data.Request.RequestURI, "/post/view/"):
		showPost(data)
	case strings.Compare(data.Request.RequestURI, "/post/comment") == 0:
		showPost(data)
	case strings.Compare(data.Request.RequestURI, "/user/login") == 0:
		userLogin(data)
	case strings.Compare(data.Request.RequestURI, "/user/register") == 0:
		userRegister(data)
	case strings.Compare(data.Request.RequestURI, "/user/logout") == 0:
		if data.Request.Method == http.MethodGet {
			userLogout(data)
		} else {
			(&models.Error{}).Consume(models.ErrorMethodNotAllowed).LogAndRespondError(data.Response, data.User)
		}
	case strings.Compare(data.Request.RequestURI, "/user/posts") == 0:
		if data.Request.Method == http.MethodGet {
			showUserPosts(data)
		} else {
			(&models.Error{}).Consume(models.ErrorMethodNotAllowed).LogAndRespondError(data.Response, data.User)
		}
	case strings.Compare(data.Request.RequestURI, "/user/likes") == 0:
		if data.Request.Method == http.MethodGet {
			showUserLikedPosts(data)
		} else {
			(&models.Error{}).Consume(models.ErrorMethodNotAllowed).LogAndRespondError(data.Response, data.User)
		}
	case strings.Compare(data.Request.RequestURI, "/user") == 0:
		if data.Request.Method == http.MethodGet {
			showUserView(data)
		} else {
			(&models.Error{}).Consume(models.ErrorMethodNotAllowed).LogAndRespondError(data.Response, data.User)
		}
	case strings.HasPrefix(data.Request.RequestURI, "/uploads/"):
		handleImages(data)
	case strings.Compare(data.Request.RequestURI, "/") == 0:
		if data.Request.Method == http.MethodGet {
			Index(data)
		} else {
			(&models.Error{}).Consume(models.ErrorMethodNotAllowed).LogAndRespondError(data.Response, data.User)
		}
	default:
		(&models.Error{}).Consume(models.ErrorNotFound).LogAndRespondError(data.Response, data.User)
	}
}

func RoutesHandler(res http.ResponseWriter, req *http.Request) {
	log.Printf("Info: %s -> %s http://%s%s", req.RemoteAddr, req.Method, req.Host, req.RequestURI)
	log.Printf("Cookies: %d", len(req.Cookies()))
	var err error
	var user models.User = models.GetGuestUser()
	for _, cookie := range req.Cookies() {
		if cookie.Name == "__Host-FRMSessionID" {
			user, err = models.GetUserBySession(cookie.Value)
			if err != nil {
				(&models.Error{}).Consume(err).LogError()
				break
			}
			user.LoggedIn = true
		}
	}
	data := models.ResponseStruct{}
	data.Init().SetResponse(res).SetRequest(req).SetUser(user)
	err = data.Request.ParseForm()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	utils.LogDebug(data.Request.Form)
	if req.Method != http.MethodPost && req.Method != http.MethodGet {
		(&models.Error{}).Consume(models.ErrorMethodNotAllowed).LogAndRespondError(data.Response, data.User)
	}
	if data.User.LoggedIn {
		notifications, err := data.User.GetNotifications()
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
			return
		}
		data.User.Notifications = notifications
	}
	Routes(data)
}
