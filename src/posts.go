package forum

import (
	"errors"
	"net/http"
	"strconv"
	"time"
)

type Post struct {
	Id        int
	Title     string
	Body      string
	UserId    int
	Timestamp time.Time
	Likes     int
	Liked     bool
	Dislikes  int
	Disliked  bool
	Category  Category
	Comments  Comments
}

type Posts []Post

func (p *Post) validatePost() error {
	if len(p.Title) == 0 {
		return ErrorPostTitleEmpty
	}
	if len(p.Body) == 0 {
		return ErrorPostBodyEmpty
	}
	if p.Category.IsEmpty() {
		return ErrorPostHasNoCategory
	}
	return nil
}

func returnMockPost(post_id int) Posts {
	return Posts{
		{
			Id:        post_id,
			Title:     "something",
			Body:      "mpla mpla",
			Timestamp: time.Now().UTC(),
			Likes:     2,
			Dislikes:  1,
			Category: Category{
				Id:   1,
				Name: "various",
			},
			Comments: ReturnMockComments(),
		},
	}
}

func showPost(res http.ResponseWriter, req *http.Request, user User) {
	data := ReturnMockResponse()
	data.User = user
	id := req.URL.Query().Get("id")
	if len(id) == 0 {
		var e Error
		e = e.Consume(ErrorPostEmptyId)
		e.LogError()
		e.RespondError(res)
		return
	}
	id_int, err := strconv.Atoi(id)
	if err != nil {
		var e Error
		e = e.Consume(err)
		e.LogError()
		e.RespondError(res)
		return
	}
	post, err := getPostById(id_int)
	if err != nil {
		var e Error
		e = e.Consume(err)
		e.LogError()
		e.RespondError(res)
		return
	}
	data.Posts = Posts{post}
	// data.Init()
	// data.SetUser(user)
	// if err != nil {
	// 	e := &Error{}
	// 	errC := e.Consume(err)
	// 	errC.LogError()
	// 	data.SetView("error_view")
	// 	data.SetError(errC)
	// 	data.WriteResponse(res)
	// 	return
	// }
	data.SetView("post_view")
	// data.SetPosts(returnMockPost(post_id))
	data.WriteResponse(res)
}

func showPosts(res http.ResponseWriter, req *http.Request, user User) {
	// query := req.URL.Query()
	// _, ok := query["id"]
	// if ok {
	// 	log.Printf("%v", query["id"])
	// 	showPost(res, query["id"][0], user)
	// 	return
	// }
	data := ReturnMockResponse()
	posts, err := getAllPosts()
	if err != nil {
		data.Error = (&Error{}).Consume(err)
		respondView(res, "error_view", data)
		return
	}
	data.SetPosts(posts)
	data.SetUser(user)
	data.SetView("posts_view")
	data.WriteResponse(res)
}

func createPostView(res http.ResponseWriter, req *http.Request, user User) {
	data := ReturnMockResponse()
	data.SetUser(user)
	data.SetView("post_create_view")
	data.WriteResponse(res)
}

func createPost(res http.ResponseWriter, req *http.Request, user User) {
	data := ReturnMockResponse()
	data.User = user

	// Parse form data
	err := req.ParseForm()
	if err != nil {
		data.Error = (&Error{}).Consume(err)
		respondView(res, "user_register_view", data)
		return
	}

	// Get form values
	title := req.FormValue("title")
	body := req.FormValue("body")
	categoryId, err := strconv.Atoi(req.FormValue("category"))
	if err != nil {
		data.Error = (&Error{}).Consume(ErrorPostHasNoCategory)
		respondView(res, "post_create_view", data)
		return
	}

	// Validate user is logged in
	if !user.LoggedIn {
		data.Error = (&Error{}).Consume(errors.New("You must be logged in to create a post"))
		respondView(res, "user_login_view", data)
		return
	}

	// Create post object
	post := Post{
		Title:     title,
		Body:      body,
		UserId:    user.Id,
		Timestamp: time.Now().UTC(),
		Category: Category{
			Id: categoryId,
		},
	}

	// Validate post
	err = post.validatePost()
	if err != nil {
		data.Error = (&Error{}).Consume(err)
		respondView(res, "post_create_view", data)
		return
	}

	// Save post to database
	postId, err := addPost(post)
	if err != nil {
		data.Error = (&Error{}).Consume(err)
		respondView(res, "post_create_view", data)
		return
	}

	// Redirect to the posts page
	http.Redirect(res, req, "/posts?id="+strconv.Itoa(postId), http.StatusSeeOther) // need to convert postId to string
}
