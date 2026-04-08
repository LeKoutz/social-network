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
	Request     *http.Request
	Response    http.ResponseWriter
	Message     Message
	Version     string
}

type ResponseStruct4ViewsIface interface {
	SetView(string) *ResponseStruct
	WriteResponse()
}

func (r ResponseStruct) WriteResponse() {
	if r.Error.StatusCode != 0 {
		r.Response.WriteHeader(r.Error.StatusCode)
	}
	respondView(r)
}

func (r *ResponseStruct) Init() *ResponseStruct {
	r.WebsiteName = forum.WebsiteName
	r.Version = getVersion()
	return r
}

func (r *ResponseStruct) SetWebsiteName(websiteName string) *ResponseStruct {
	r.WebsiteName = websiteName
	return r
}

func (r ResponseStruct) SetView(viewname string) *ResponseStruct {
	r.View = viewname
	return &r
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

func (r *ResponseStruct) SetErrorConsume(err error) *ResponseStruct {
	r.Error.Consume(err)
	return r
}

func (r *ResponseStruct) SetRequest(req *http.Request) *ResponseStruct {
	r.Request = req
	return r
}

func (r *ResponseStruct) SetResponse(res http.ResponseWriter) *ResponseStruct {
	r.Response = res
	return r
}

func (r *ResponseStruct) GetResponse(res http.ResponseWriter) http.ResponseWriter {
	return r.Response
}

func respondView(data ResponseStruct) {
	var templatesDir string = "templates"
	var tmpl *template.Template
	tmpl, err := template.ParseGlob(templatesDir + "/*.html")
	if err != nil {
		(&Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	err = tmpl.ExecuteTemplate(data.Response, data.View, data)
	if err != nil {
		(&Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
}
