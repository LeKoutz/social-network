package controllers

import (
	"context"
	"forum/src/models"
	"forum/src/utils"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func RegisterRoutes(mux *http.ServeMux) {
	// mux.HandleFunc("GET /user/login", userLogin)
	// mux.HandleFunc("POST /user/login", showPost)
	// mux.HandleFunc("GET /user/logout", showPost)
	// mux.HandleFunc("POST /user/logout", showPost)
	// mux.HandleFunc("GET /user/register", showPost)
	// mux.HandleFunc("POST /user/register", showPost)
	// mux.HandleFunc("GET /post/{id}/view", showPost)
	// mux.HandleFunc("GET /post/new", showPost)
	// mux.HandleFunc("POST /post/new", showPost)
	// mux.HandleFunc("POST /post/{id}/react", showPost)
	// mux.HandleFunc("POST /post/{id}/comment", showPost)
	// mux.HandleFunc("GET /category/{id}, showCategory)
	// mux.HandleFunc("GET /", Index)
}

func Routes(data models.ResponseStruct) {
	uri, err := url.ParseRequestURI(data.Request.RequestURI)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	utils.LogDebug(uri)

	categoryMask, err := regexp.Compile("/category/[0-9]+")
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	postMask, err := regexp.Compile("/post/[0-9]+")
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	switch {
	case categoryMask.MatchString(data.Request.RequestURI):
		showCategory(data)
	case postMask.MatchString(data.Request.RequestURI):
		handlePost(data)
	case strings.HasPrefix(data.Request.RequestURI, "/comment"):
		handleComment(data)
	case strings.Compare(data.Request.RequestURI, "/categories") == 0:
		showCategories(data)
	case strings.Compare(data.Request.RequestURI, "/posts") == 0:
		showPosts(data)
	case strings.Compare(data.Request.RequestURI, "/post/create") == 0:
		handlePost(data)
	case strings.HasPrefix(data.Request.RequestURI, "/post/view/"):
		showPost(data)
	case strings.Compare(data.Request.RequestURI, "/post/comment") == 0:
		showPost(data)
	case strings.Compare(data.Request.RequestURI, "/user/login") == 0:
		userLogin(data)
	case strings.Compare(data.Request.RequestURI, "/user/register") == 0:
		userRegister(data)
	case strings.Compare(data.Request.RequestURI, "/user/logout") == 0:
		userLogout(data)
	case strings.Compare(data.Request.RequestURI, "/user/posts") == 0:
		showUserPosts(data)
	case strings.Compare(data.Request.RequestURI, "/user/likes") == 0:
		showUserLikedPosts(data)
	case strings.Compare(data.Request.RequestURI, "/") == 0:
		Index(data)
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
	Routes(data)
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		log.Printf("Info: %s -> %s http://%s%s", req.RemoteAddr, req.Method, req.Host, req.RequestURI)
		log.Printf("Cookies: %d", len(req.Cookies()))
		next.ServeHTTP(res, req)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
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
		// irrelevant anymore I guess
		// data := models.ResponseStruct{}
		// data.Init().SetResponse(res).SetRequest(req).SetUser(user)
		// err = data.Request.ParseForm()
		// if err != nil {
		// 	(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		// 	return
		// }
		// utils.LogDebug(data.Request.Form)
		// if req.Method != http.MethodPost && req.Method != http.MethodGet {
		// 	(&models.Error{}).Consume(models.ErrorMethodNotAllowed).LogAndRespondError(data.Response, data.User)
		// }
		ctx := context.WithValue(req.Context(), "User", user)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}
