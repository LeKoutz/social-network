package controllers

import (
	"forum/src/models"
	"forum/src/views"
	"log"
	"net/http"
	"strings"
)

func getRoutes(data models.ResponseStruct) {
	switch {
	case strings.HasPrefix(data.Request.RequestURI, "/posts"):
		showPosts(data)
	case strings.HasPrefix(data.Request.RequestURI, "/post?action=new"):
		views.PostCreate(data)
	case strings.HasPrefix(data.Request.RequestURI, "/post"):
		showPost(data)
	case strings.HasPrefix(data.Request.RequestURI, "/login"):
		userLogin(data)
	case strings.HasPrefix(data.Request.RequestURI, "/register"):
		views.UserRegister(data)
	case strings.HasPrefix(data.Request.RequestURI, "/logout"):
		userLogout(data)
	case strings.HasPrefix(data.Request.RequestURI, "/categories"):
		showCategories(data)
	case strings.HasPrefix(data.Request.RequestURI, "/category?id="):
		showCategory(data)
	case strings.HasPrefix(data.Request.RequestURI, "/my/posts"):
		showUserPosts(data)
	case strings.HasPrefix(data.Request.RequestURI, "/my/liked"):
		showUserLikedPosts(data)
	case strings.Compare(data.Request.RequestURI, "/") == 0:
		Index(data)
	default:
		(&models.Error{}).Consume(models.ErrorNotFound).LogAndRespondError(data.Response, data.User)
	}
}

func postRoutes(data models.ResponseStruct) {
	switch {
	case strings.HasPrefix(data.Request.RequestURI, "/post?action=create"):
		createPost(data)
	case strings.HasPrefix(data.Request.RequestURI, "/user?action=login"):
		attemptLogin(data)
	case strings.HasPrefix(data.Request.RequestURI, "/user?action=register"):
		registerUser(data)
	case strings.HasPrefix(data.Request.RequestURI, "/comment?action=create&post_id="):
		createComment(data)
	case strings.HasPrefix(data.Request.RequestURI, "/post"):
		handlePostReaction(data)
	case strings.HasPrefix(data.Request.RequestURI, "/comment"):
		handleCommentReaction(data)
	case strings.Compare(data.Request.RequestURI, "/categories") == 0:
		showCategories(data)
	case strings.Compare(data.Request.RequestURI, "/") == 0:
		(&models.Error{}).Consume(models.ErrorMethodNotAllowed).LogAndRespondError(data.Response, data.User)
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
	switch req.Method {
	case http.MethodGet:
		getRoutes(data)
	case http.MethodPost:
		postRoutes(data)
	default:
		(&models.Error{}).Consume(models.ErrorMethodNotAllowed).LogAndRespondError(data.Response, data.User)
	}
}
