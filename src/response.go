package forum

import "net/http"

type ResponseStruct struct {
	WebsiteName string
	View        string
	User        User
	Posts       Posts
	Categories  Categories
	Error       Error
}

func (r *ResponseStruct) WriteResponse(res http.ResponseWriter) {
	respondView(res, r.View, *r)
}

func (r *ResponseStruct) SetWebsiteName(websiteName string) {
	r.WebsiteName = websiteName
}

func (r *ResponseStruct) SetView(viewname string) {
	r.View = viewname
}

func (r *ResponseStruct) SetUser(user User) {
	r.User = user
}

func (r *ResponseStruct) SetPosts(posts Posts) {
	r.Posts = posts
}

func (r *ResponseStruct) SetCategories(categories Categories) {
	r.Categories = categories
}

func (r *ResponseStruct) SetError(err Error) {
	r.Error = err
}
