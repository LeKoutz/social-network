package controllers

import (
	"forum/src/models"
	"forum/src/utils"
	"log"
	"net/http"
	"strings"
)

type Routes []Route

type RouteController func(models.ResponseStruct)

type Route struct {
	Function RouteController
	Prefix   bool
	Method   string
	Path     string
	NeedsLogin bool
}

var routes = Routes{
	Route{Method: "GET", Path: "/api/", Function: Index},

	Route{Method: "GET", Path: "/api/category/view/", Prefix: true, Function: showCategory, NeedsLogin: true},
	Route{Method: "GET", Path: "/api/categories", Prefix: true, Function: showCategories, NeedsLogin: true},

	Route{Method: "POST", Path: "/api/comment/create", Function: handleCommentCreate, NeedsLogin: true},
	Route{Method: "POST", Path: "/api/comment/react", Function: handleCommentReaction, NeedsLogin: true},
	Route{Method: "POST", Path: "/api/comment/edit", Function: handleCommentEdit, NeedsLogin: true},
	Route{Method: "POST", Path: "/api/comment/delete", Function: handleCommentDelete, NeedsLogin: true},

	Route{Method: "GET", Path: "/api/auth/google/callback", Prefix: true, Function: handleGoogleCallback},
	Route{Method: "GET", Path: "/api/auth/google", Prefix: true, Function: handleOAuthLoginGoogle},
	Route{Method: "GET", Path: "/api/auth/github/callback", Prefix: true, Function: handleGitHubCallback},
	Route{Method: "GET", Path: "/api/auth/github", Prefix: true, Function: handleOAuthLoginGithub},

	Route{Method: "GET", Path: "/api/posts", Function: showPosts, NeedsLogin: true},
	Route{Method: "*", Path: "/api/post/create", Function: handlePostCreate, NeedsLogin: true},
	Route{Method: "POST", Path: "/api/post/react", Function: handlePostReaction, NeedsLogin: true},
	Route{Method: "GET", Path: "/api/post/view/", Prefix: true, Function: showPost, NeedsLogin: true},
	Route{Method: "GET", Path: "/api/post/comment", Function: showPost, NeedsLogin: true},
	Route{Method: "*", Path: "/api/post/edit", Prefix: true, Function: handlePostEdit, NeedsLogin: true},
	Route{Method: "POST", Path: "/api/post/delete", Function: handlePostDelete, NeedsLogin: true},

	Route{Method: "*", Path: "/api/user/login", Function: userLogin},
	Route{Method: "*", Path: "/api/user/register", Function: userRegister},
	Route{Method: "GET", Path: "/api/user/logout", Function: userLogout, NeedsLogin: true},
	Route{Method: "GET", Path: "/api/user/posts", Function: showUserPosts, NeedsLogin: true},
	Route{Method: "GET", Path: "/api/user/likes", Function: showUserLikedPosts, NeedsLogin: true},
	Route{Method: "POST", Path: "/api/user/notifications", Function: markAllNotificationsAsRead, NeedsLogin: true},
	Route{Method: "GET", Path: "/api/user", Function: showUserView, NeedsLogin: true},
	Route{Method: "GET", Path: "/api/user/activity", Function: showUserActivity, NeedsLogin: true},

	Route{Method: "GET", Path: "/ws", Function: serveWs, NeedsLogin: true},
	Route{Method: "GET", Path: "/api/users", Function: getUsers, NeedsLogin: true},
	Route{Method: "GET", Path: "/api/chat/", Prefix: true, Function: showChatHistory, NeedsLogin: true},

	Route{Method: "GET", Path: "/uploads/", Prefix: true, Function: handleImages, NeedsLogin: true},
	Route{Method: "GET", Path: "/", Prefix: true, Function: serveSPA},
}

func matchRoute(data models.ResponseStruct) (*Route, error) {
	for _, route_s := range routes {
		if route_s.Prefix && strings.HasPrefix(data.Request.RequestURI, route_s.Path) {
			if route_s.Method == data.Request.Method || route_s.Method == "*" {
				return &route_s, nil
			} else {
				return nil, models.ErrorMethodNotAllowed
			}
		} else if strings.Compare(data.Request.RequestURI, route_s.Path) == 0 {
			if route_s.Method == data.Request.Method || route_s.Method == "*" {
				return &route_s, nil
			} else {
				return nil, models.ErrorMethodNotAllowed
			}
		}
	}
	return nil, models.ErrorNotFound
}

func RouteToController(data models.ResponseStruct) {
	route, err := matchRoute(data)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	if route != nil {
		if data.User.LoggedIn && route.NeedsLogin {
			route.Function(data)
		} else if !route.NeedsLogin {
			route.Function(data)
		} else {
			(&models.Error{}).Consume(models.ErrorUnauthorizedAction).LogAndRespondError(data.Response, data.User)
		}
		return
	}
	(&models.Error{}).Consume(models.ErrorNotFound).LogAndRespondError(data.Response, data.User)
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
			err := user.GetNotifications()
			if err != nil {
				(&models.Error{}).Consume(err).LogError()
			}
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
	RouteToController(data)
}
