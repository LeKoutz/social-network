package forum

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

func getRoutes(res http.ResponseWriter, req *http.Request, user User) {
	switch {
	case strings.HasPrefix(req.RequestURI, "/posts"):
		showPosts(res, req, user)
	case strings.HasPrefix(req.RequestURI, "/post?action=new"):
		createPostView(res, req, user)
	case strings.HasPrefix(req.RequestURI, "/post"):
		showPost(res, req, user)
	case strings.HasPrefix(req.RequestURI, "/login"):
		showLogin(res, req, user)
	case strings.HasPrefix(req.RequestURI, "/register"):
		showRegister(res, req, user)
	case strings.HasPrefix(req.RequestURI, "/logout"):
		showLogout(res, req, user)
	case strings.HasPrefix(req.RequestURI, "/categories"):
		showCategories(res, req, user)
	case strings.HasPrefix(req.RequestURI, "/category?id="):
		showCategory(res, req, user)
	case strings.HasPrefix(req.RequestURI, "/my"):
		showUserStuff(res, req, user)
	case strings.Compare(req.RequestURI, "/") == 0:
		showIndex(res, req, user)
	default:
		res.WriteHeader(http.StatusNotFound)
		(&Error{}).Consume(ErrorNotFound).LogAndRespondError(res, user)
	}
}

func postRoutes(res http.ResponseWriter, req *http.Request, user User) {
	switch {
	case strings.HasPrefix(req.RequestURI, "/post?action=create"):
		createPost(res, req, user)
	case strings.HasPrefix(req.RequestURI, "/user?action=login"):
		attemptLogin(res, req, user)
	case strings.HasPrefix(req.RequestURI, "/user?action=register"):
		registerUser(res, req)
	case strings.HasPrefix(req.RequestURI, "/comment?action=create&post_id="):
		createComment(res, req, user)
	case strings.HasPrefix(req.RequestURI, "/post"):
		handlePostReaction(res, req, user)
	case strings.HasPrefix(req.RequestURI, "/comment"):
		handleCommentReaction(res, req, user)
	case strings.Compare(req.RequestURI, "/categories") == 0:
		showCategories(res, req, user)
	case strings.Compare(req.RequestURI, "/") == 0:
		showIndex(res, req, user)
	default:
		res.WriteHeader(http.StatusNotFound)
		(&Error{}).Consume(ErrorNotFound).LogAndRespondError(res, user)
	}
}

func routesHandler(res http.ResponseWriter, req *http.Request) {
	log.Printf("Info: %s -> %s http://%s%s", req.RemoteAddr, req.Method, req.Host, req.RequestURI)
	log.Printf("Cookies: %d", len(req.Cookies()))
	var err error
	var user User = GuestUser
	for _, cookie := range req.Cookies() {
		if cookie.Name == "__Host-FRMSessionID" {
			user, err = getUserBySession(cookie.Value)
			if err != nil {
				(&Error{}).Consume(err).LogError()
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
		res.WriteHeader(http.StatusMethodNotAllowed)
		(&Error{}).Consume(ErrorNotFound).LogAndRespondError(res, user)
	}
}

func respondError(statusInt int, res http.ResponseWriter, _ string) {
	res.WriteHeader(statusInt)
	var templatesDir string = "templates"
	var index *template.Template
	index, err := template.ParseGlob(templatesDir + "/*.html")
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return
	}
	err = index.ExecuteTemplate(res, "error_view", ReturnMockResponse())
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return
	}
}
