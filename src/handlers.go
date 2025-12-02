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
	case strings.HasPrefix(req.RequestURI, "/login"):
		showLogin(res, req, user)
	case strings.HasPrefix(req.RequestURI, "/register"):
		showRegister(res, req, user)
	case strings.HasPrefix(req.RequestURI, "/logout"):
		showLogout(res, req, user)
	case strings.HasPrefix(req.RequestURI, "/categories"):
		showCategories(res, req, user)
	case strings.Compare(req.RequestURI, "/") == 0:
		showIndex(res, req, user)
	default:
		log.Printf("%s", req.RequestURI)
		res.WriteHeader(http.StatusNotFound)
		respondView(res, "error_view", ResponseStruct{
			WebsiteName: "Forum",
			Error: Error{
				True:    true,
				Message: "Not found",
			},
		})
	}
}

func postRoutes(res http.ResponseWriter, req *http.Request, user User) {
	switch {
	case strings.HasPrefix(req.RequestURI, "/posts"):
		showPosts(res, req, user)
	case strings.HasPrefix(req.RequestURI, "/user?action=login"):
		attemptLogin(res, req, user)
	case strings.HasPrefix(req.RequestURI, "/user?action=register"):
		showRegister(res, req, user)
	case strings.HasPrefix(req.RequestURI, "/categories"):
		showCategories(res, req, user)
	case strings.Compare(req.RequestURI, "/") == 0:
		showIndex(res, req, user)
	default:
		log.Printf("%s", req.RequestURI)
		res.WriteHeader(http.StatusNotFound)
		respondView(res, "error_view", ResponseStruct{
			WebsiteName: "Forum",
			Error: Error{
				True:    true,
				Message: "Not found",
			},
		})
	}
}

func routesHandler(res http.ResponseWriter, req *http.Request) {
	log.Printf("Info: %s -> %s http://%s%s", req.RemoteAddr, req.Method, req.Host, req.RequestURI)
	log.Printf("Cookies: %d", len(req.Cookies()))
	var user User = User{
		Name: "Guest",
		Email: "guest@example.com",
	}
	for _, cookie := range req.Cookies() {
		log.Printf("%#v", cookie)
		if cookie.Name == "access" && cookie.Value == "admin" {
			user = AdminUser
		}
	}
	switch req.Method {
	case http.MethodGet:
		getRoutes(res, req, user)
	case http.MethodPost:
		postRoutes(res, req, user)
	default:
		res.WriteHeader(http.StatusMethodNotAllowed)
		respondView(res, "error_view", ResponseStruct{
			WebsiteName: "Forum",
			Error: Error{
				True:    true,
				Message: "Method Not Allowed",
			},
		})
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
	err = index.ExecuteTemplate(res, "error_view", ValuesToClient())
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return
	}
}
