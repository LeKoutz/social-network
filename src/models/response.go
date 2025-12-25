package models

import (
	"forum/src"
	"html/template"
	"net/http"
)

type ResponseStruct struct {
	WebsiteName string
	View        string
	User        User
	Posts       Posts
	Categories  Categories
	Error       Error
}

func (r *ResponseStruct) WriteResponse(res http.ResponseWriter) {
	if r.Error.StatusCode != 0 {
		res.WriteHeader(r.Error.StatusCode)
	}
	respondView(res, *r)
}

func (r *ResponseStruct) Init() *ResponseStruct {
	r.WebsiteName = forum.WebsiteName
	return r
}

func (r *ResponseStruct) SetWebsiteName(websiteName string) *ResponseStruct {
	r.WebsiteName = websiteName
	return r
}

func (r *ResponseStruct) SetView(viewname string) *ResponseStruct {
	r.View = viewname
	return r
}

func (r *ResponseStruct) SetUser(user User) *ResponseStruct {
	r.User = user
	return r
}

func (r *ResponseStruct) SetPosts(posts Posts) *ResponseStruct {
	r.Posts = posts
	return r
}

func (r *ResponseStruct) SetCategories(categories Categories) *ResponseStruct {
	r.Categories = categories
	return r
}

func (r *ResponseStruct) SetError(err Error) *ResponseStruct {
	r.Error = err
	return r
}

func respondView(res http.ResponseWriter, data ResponseStruct) {
	var templatesDir string = "templates"
	var tmpl *template.Template
	tmpl, err := template.ParseGlob(templatesDir + "/*.html")
	if err != nil {
		(&Error{}).Consume(err).LogAndRespondError(res, data.User)
		return
	}
	err = tmpl.ExecuteTemplate(res, data.View, data)
	if err != nil {
		(&Error{}).Consume(err).LogAndRespondError(res, data.User)
		return
	}
}
