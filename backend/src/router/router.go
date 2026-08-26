package router

import (
	"forum/src/ferror"
	"forum/src/handlers"
	"forum/src/middleware"
	"forum/src/state"
	"log"
	"net/http"
	"strings"
)

type Routes []Route

type RouteController func(state.StateHandler)

type Route struct {
	Function   RouteController
	Prefix     bool
	Method     string
	Path       string
	NeedsLogin bool
}

var routes = Routes{
	Route{Method: "GET", Path: "/api", Function: handlers.HandleIndex},

	// Categories
	Route{Method: "GET", Path: "/api/category/view/", Prefix: true, Function: handlers.HandleShowCategory, NeedsLogin: true},
	Route{Method: "GET", Path: "/api/categories", Prefix: true, Function: handlers.HandleShowCategories, NeedsLogin: true},

	// Comments
	Route{Method: "POST", Path: "/api/comment/create", Function: handlers.HandleCommentCreate, NeedsLogin: true},
	Route{Method: "POST", Path: "/api/comment/react", Function: handlers.HandleCommentReaction, NeedsLogin: true},
	Route{Method: "POST", Path: "/api/comment/edit", Function: handlers.HandleCommentEdit, NeedsLogin: true},
	Route{Method: "POST", Path: "/api/comment/delete", Function: handlers.HandleCommentDelete, NeedsLogin: true},

	// Auth
	Route{Method: "GET", Path: "/api/auth/google/callback", Prefix: true, Function: handlers.HandleGoogleCallback},
	Route{Method: "GET", Path: "/api/auth/google", Prefix: true, Function: handlers.HandleOAuthLoginGoogle},
	Route{Method: "GET", Path: "/api/auth/github/callback", Prefix: true, Function: handlers.HandleGitHubCallback},
	Route{Method: "GET", Path: "/api/auth/github", Prefix: true, Function: handlers.HandleOAuthLoginGithub},

	// Post(s)
	Route{Method: "GET", Path: "/api/posts", Function: handlers.HandleShowPosts, NeedsLogin: true},
	Route{Method: "*", Path: "/api/post/create", Function: handlers.HandlePostCreate, NeedsLogin: true},
	Route{Method: "POST", Path: "/api/post/react", Function: handlers.HandlePostReaction, NeedsLogin: true},
	Route{Method: "GET", Path: "/api/post/view/", Prefix: true, Function: handlers.HandleShowPost, NeedsLogin: true},
	Route{Method: "GET", Path: "/api/post/comment", Function: handlers.HandleShowPost, NeedsLogin: true},
	Route{Method: "*", Path: "/api/post/edit", Prefix: true, Function: handlers.HandlePostEdit, NeedsLogin: true},
	Route{Method: "POST", Path: "/api/post/delete", Function: handlers.HandlePostDelete, NeedsLogin: true},

	// User
	Route{Method: "*", Path: "/api/user/login", Function: handlers.HandleUserLogin},
	Route{Method: "*", Path: "/api/user/register", Function: handlers.HandleUserRegister},
	Route{Method: "GET", Path: "/api/user/logout", Function: handlers.HandleUserLogout, NeedsLogin: true},
	Route{Method: "GET", Path: "/api/user/posts", Function: handlers.HandleShowUserPosts, NeedsLogin: true},
	Route{Method: "GET", Path: "/api/user/likes", Function: handlers.HandleShowUserLikedPosts, NeedsLogin: true},
	Route{Method: "POST", Path: "/api/user/notifications", Function: handlers.HandleMarkAllNotificationsAsRead, NeedsLogin: true},
	Route{Method: "GET", Path: "/api/user", Function: handlers.HandleShowUserView, NeedsLogin: true},
	Route{Method: "GET", Path: "/api/user/activity", Function: handlers.HandleShowUserActivity, NeedsLogin: true},

	// WebSocket
	Route{Method: "GET", Path: "/ws", Function: handlers.HandleWs, NeedsLogin: true},

	// Chat
	Route{Method: "GET", Path: "/api/users", Function: handlers.HandleGetUsers, NeedsLogin: true},
	Route{Method: "GET", Path: "/api/chat/unread", Function: handlers.HandleServeUnreadMessages, NeedsLogin: true},
	Route{Method: "GET", Path: "/api/chat/", Prefix: true, Function: handlers.HandleShowChatHistory, NeedsLogin: true},

	// Uploads
	Route{Method: "GET", Path: "/uploads/", Prefix: true, Function: handlers.HandleImages, NeedsLogin: true},

	// SPA
	// TODO Remove
	Route{Method: "GET", Path: "/", Prefix: true, Function: handlers.HandleServeSPA},
}

func matchRoute(data state.StateRoute) (*Route, error) {
	for _, route_s := range routes {
		if route_s.Prefix && strings.HasPrefix(data.GetRequest().RequestURI, route_s.Path) {
			if route_s.Method == data.GetRequest().Method || route_s.Method == "*" {
				return &route_s, nil
			} else {
				return nil, ferror.ErrorMethodNotAllowed
			}
		} else if strings.Compare(data.GetRequest().RequestURI, route_s.Path) == 0 {
			if route_s.Method == data.GetRequest().Method || route_s.Method == "*" {
				return &route_s, nil
			} else {
				return nil, ferror.ErrorMethodNotAllowed
			}
		}
	}
	return nil, ferror.ErrorNotFound
}

func RouteToController(data state.StateRoute) {
	route, err := matchRoute(data)
	if err != nil {
		data.SetErrorConsume(err)
		data.WriteResponse()
		return
	}
	if route != nil {
		if data.GetUser().LoggedIn && route.NeedsLogin {
			route.Function(data.(state.StateHandler))
		} else if !route.NeedsLogin {
			route.Function(data.(state.StateHandler))
		} else {
			data.SetErrorConsume(ferror.ErrorUnauthorizedAction)
			data.WriteResponse()
		}
		return
	}
	data.SetErrorConsume(ferror.ErrorNotFound)
	data.WriteResponse()
}

func RoutesHandler(res http.ResponseWriter, req *http.Request) {
	log.Printf("Info: %s -> %s http://%s%s", req.RemoteAddr, req.Method, req.Host, req.RequestURI)
	log.Printf("Cookies: %d", len(req.Cookies()))
	var err error
	data := state.State{}
	data.Init().SetResponse(res).SetRequest(req)
	middleware.AuthMiddleware(&data)
	err = data.Request.ParseForm()
	if err != nil {
		data.SetErrorConsume(err)
		data.WriteResponse()
		return
	}
	RouteToController(&data)
}
