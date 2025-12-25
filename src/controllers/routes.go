package controllers

import (
	"forum/src/models"
	"forum/src/views"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func getRoutes(res http.ResponseWriter, req *http.Request, user models.User) {
	switch {
	case strings.HasPrefix(req.RequestURI, "/posts"):
		showPosts(res, req, user)
	case strings.HasPrefix(req.RequestURI, "/post?action=new"):
		views.PostCreate(res, req, user)
	case strings.HasPrefix(req.RequestURI, "/post"):
		showPost(res, req, user)
	case strings.HasPrefix(req.RequestURI, "/login"):
		views.UserLogin(res, req, user)
	case strings.HasPrefix(req.RequestURI, "/register"):
		views.UserRegister(res, req, user)
	case strings.HasPrefix(req.RequestURI, "/logout"):
		views.UserLogout(res, req, user)
	case strings.HasPrefix(req.RequestURI, "/categories"):
		showCategories(res, req, user)
	case strings.HasPrefix(req.RequestURI, "/category?id="):
		showCategory(res, req, user)
	case strings.HasPrefix(req.RequestURI, "/my/posts"):
		showUserPosts(res, user)
	case strings.HasPrefix(req.RequestURI, "/my/liked"):
		showUserLikedPosts(res, user)
	case strings.Compare(req.RequestURI, "/") == 0:
		Index(res, req, user)
	default:
		(&models.Error{}).Consume(models.ErrorNotFound).LogAndRespondError(res, user)
	}
}

func postRoutes(res http.ResponseWriter, req *http.Request, user models.User) {
	switch {
	case strings.HasPrefix(req.RequestURI, "/post?action=create"):
		createPost(res, req, user)
	case strings.HasPrefix(req.RequestURI, "/user?action=login"):
		attemptLogin(res, req, user)
	case strings.HasPrefix(req.RequestURI, "/user?action=register"):
		registerUser(res, req, user)
	case strings.HasPrefix(req.RequestURI, "/comment?action=create&post_id="):
		createComment(res, req, user)
	case strings.HasPrefix(req.RequestURI, "/post"):
		handlePostReaction(res, req, user)
	case strings.HasPrefix(req.RequestURI, "/comment"):
		handleCommentReaction(res, req, user)
	case strings.Compare(req.RequestURI, "/categories") == 0:
		showCategories(res, req, user)
	case strings.Compare(req.RequestURI, "/") == 0:
		(&models.Error{}).Consume(models.ErrorMethodNotAllowed).LogAndRespondError(res, user)
	default:
		(&models.Error{}).Consume(models.ErrorNotFound).LogAndRespondError(res, user)
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
	switch req.Method {
	case http.MethodGet:
		getRoutes(res, req, user)
	case http.MethodPost:
		postRoutes(res, req, user)
	default:
		(&models.Error{}).Consume(models.ErrorMethodNotAllowed).LogAndRespondError(res, user)
	}
}

func respondError(statusInt int, res http.ResponseWriter, _ string) {
	var templatesDir string = "templates"
	var index *template.Template
	index, err := template.ParseGlob(templatesDir + "/*.html")
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return
	}
	err = index.ExecuteTemplate(res, "error_view", (&models.ResponseStruct{}).Init())
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return
	}
}
