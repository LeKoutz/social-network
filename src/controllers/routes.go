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
}

var routes = Routes{
	Route{Method: "GET", Path: "/", Function: Index},

	Route{Method: "GET", Path: "/category/view/", Prefix: true, Function: showCategory},
	Route{Method: "GET", Path: "/categories", Prefix: true, Function: showCategories},

	Route{Method: "POST", Path: "/comment/create", Function: handleCommentCreate},
	Route{Method: "POST", Path: "/comment/react", Function: handleCommentReaction},
	Route{Method: "POST", Path: "/comment/edit", Function: handleCommentEdit},
	Route{Method: "POST", Path: "/comment/delete", Function: handleCommentDelete},

	Route{Method: "GET", Path: "/auth/google/callback", Prefix: true, Function: handleGoogleCallback},
	Route{Method: "GET", Path: "/auth/google", Prefix: true, Function: handleOAuthLoginGoogle},
	Route{Method: "GET", Path: "/auth/github/callback", Prefix: true, Function: handleGitHubCallback},
	Route{Method: "GET", Path: "/auth/github", Prefix: true, Function: handleOAuthLoginGithub},

	Route{Method: "GET", Path: "/posts", Function: showPosts},
	Route{Method: "*", Path: "/post/create", Function: handlePostCreate},
	Route{Method: "POST", Path: "/post/react", Function: handlePostReaction},
	Route{Method: "GET", Path: "/post/view/", Prefix: true, Function: showPost},
	Route{Method: "GET", Path: "/post/comment", Function: showPost},
	Route{Method: "*", Path: "/post/edit", Prefix: true, Function: handlePostEdit},
	Route{Method: "POST", Path: "/post/delete", Function: handlePostDelete},

	Route{Method: "*", Path: "/user/login", Function: userLogin},
	Route{Method: "*", Path: "/user/register", Function: userRegister},
	Route{Method: "GET", Path: "/user/logout", Function: userLogout},
	Route{Method: "GET", Path: "/user/posts", Function: showUserPosts},
	Route{Method: "GET", Path: "/user/likes", Function: showUserLikedPosts},
	Route{Method: "GET", Path: "/user/notifications", Function: markAllNotificationsAsRead},
	Route{Method: "GET", Path: "/user", Function: showUserView},
	Route{Method: "GET", Path: "/user/activity", Function: showUserActivity},

	Route{Method: "GET", Path: "/uploads/", Prefix: true, Function: handleImages},
}

func matchRoute(data models.ResponseStruct) *Route {
	for _, route_s := range routes {
		if route_s.Prefix && strings.HasPrefix(data.Request.RequestURI, route_s.Path) {
			if route_s.Method == data.Request.Method || route_s.Method == "*" {
				return &route_s
			}
		} else if strings.Compare(data.Request.RequestURI, route_s.Path) == 0 {
			if route_s.Method == data.Request.Method || route_s.Method == "*" {
				return &route_s
			}
		}
	}
	return nil
}

func RouteToController(data models.ResponseStruct) {
	route := matchRoute(data)
	if route != nil {
		route.Function(data)
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
	RouteToController(data)
}
