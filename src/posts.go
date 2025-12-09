package forum

import (
	"net/http"
	"strconv"
	"time"
)

type Post struct {
	Id              int
	Title           string
	Body            string
	UserId          int
	User            User
	Timestamp       int64
	TimestampString string
	Likes           int
	Liked           bool
	Dislikes        int
	Disliked        bool
	Category        Category
	Categories      Categories
	Comments        Comments
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
			Timestamp: time.Now().Unix(),
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
		(&Error{}).Consume(ErrorPostEmptyId).LogAndRespondError(res, user)
		return
	}
	id_int, err := strconv.Atoi(id)
	if err != nil {
		(&Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	post, err := getPostById(id_int)
	if err != nil {
		(&Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	comments, err := getCommentsByPostId(id_int)
	if err != nil {
		(&Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	post.Comments = comments
	categories, err := getCategoriesByPostId(post.Id)
	if err != nil {
		(&Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	post.Categories = categories
	data.Posts = Posts{post}
	data.SetView("post_view")
	data.WriteResponse(res)
}

func showPosts(res http.ResponseWriter, _ *http.Request, user User) {
	data := ReturnMockResponse()
	posts, err := getAllPosts()
	if err != nil {
		data.Error = *(&Error{}).Consume(err)
		respondView(res, "error_view", data)
		return
	}
	data.SetPosts(posts)
	data.SetUser(user)
	data.SetView("posts_view")
	data.WriteResponse(res)
}

func createPostView(res http.ResponseWriter, _ *http.Request, user User) {
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
		data.Error = *(&Error{}).Consume(err)
		respondView(res, "user_register_view", data)
		return
	}

	// Get form values
	title := req.FormValue("title")
	body := req.FormValue("body")
	categoryId, err := strconv.Atoi(req.FormValue("category"))
	if err != nil {
		data.Error = *(&Error{}).Consume(err)
		respondView(res, "post_create_view", data)
		return
	}

	// Validate user is logged in
	if !user.LoggedIn {
		data.Error = *(&Error{}).Consume(ErrorPostPermissionDenied)
		respondView(res, "user_login_view", data)
		return
	}

	// Create post object
	post := Post{
		Title:  title,
		Body:   body,
		UserId: user.Id,
		Category: Category{
			Id: categoryId,
		},
	}

	// Save post to database
	postId, err := addPost(post)
	if err != nil {
		data.Error = *(&Error{}).Consume(err)
		respondView(res, "post_create_view", data)
		return
	}

	postIdStr := strconv.Itoa(postId)
	redirectURL := "/post?id=" + postIdStr
	// Redirect to the post's page
	http.Redirect(res, req, redirectURL, http.StatusSeeOther)
}
