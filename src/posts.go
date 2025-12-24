package forum

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type Post struct {
	Id              int64
	Title           string
	Body            string
	UserId          int64
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
	if p.Categories.IsEmpty() {
		return ErrorPostHasNoCategory
	}
	return nil
}

func returnMockPost(post_id int64) Posts {
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
	data := ResponseStruct{}
	data.Init().SetUser(user)
	id := req.URL.Query().Get("id")
	if len(id) == 0 {
		(&Error{}).Consume(ErrorPostEmptyId).LogAndRespondError(res, user)
		return
	}
	ok, err := regexp.MatchString(`^\d+$`, id)
	if !ok {
		(&Error{}).Consume(ErrorInvalidPostId).LogAndRespondError(res, user)
		return
	}
	id_int, err := strconv.ParseInt(id, 10, 64)
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
	for i := range comments {
		comments[i].getReactions()
		comments[i].getReactionsByUserId(user.Id)
	}
	post.Comments = comments
	categories, err := getCategoriesByPostId(post.Id)
	if err != nil {
		(&Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	post.Categories = categories
	err = post.getReactions()
	if err != nil {
		(&Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	err = post.getReactionsByUserId(user.Id)
	if err != nil {
		(&Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	data.Posts = Posts{post}
	data.SetView("post_view").WriteResponse(res)
}

func showPosts(res http.ResponseWriter, _ *http.Request, user User) {
	data := ResponseStruct{}
	data.Init()
	posts, err := getAllPosts()
	if err != nil {
		data.Error = *(&Error{}).Consume(err)
		(&Error{}).Consume(err).RespondError(res, user)
		return
	}
	for i := range posts {
		err = posts[i].getReactions()
		if err != nil {
			(&Error{}).Consume(err).LogAndRespondError(res, user)
			return
		}
		err = posts[i].getReactionsByUserId(user.Id)
		if err != nil {
			(&Error{}).Consume(err).LogAndRespondError(res, user)
			return
		}
	}
	data.SetPosts(posts).SetUser(user).SetView("posts_view").WriteResponse(res)
}

func createPostView(res http.ResponseWriter, _ *http.Request, user User) {
	data := ResponseStruct{}
	data.Init().SetUser(user)
	categories, err := getAllCategories()
	if err != nil {
		(&Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	data.SetCategories(categories)
	data.SetView("post_create_view").WriteResponse(res)
}

func createPost(res http.ResponseWriter, req *http.Request, user User) {
	data := ResponseStruct{}
	data.Init().SetUser(user)

	err := req.ParseForm()
	if err != nil {
		data.Error = *(&Error{}).Consume(err)
		data.SetView("user_register_view").WriteResponse(res)
		return
	}
	// Get form values
	title := req.FormValue("title")
	body := req.FormValue("body")
	categories, err := getAllCategories()
	if err != nil {
		data.Error = *(&Error{}).Consume(err)
		data.SetView("post_create_view").WriteResponse(res)
		return
	}
	var PostCategories Categories
	for _, category := range categories {
		cc := fmt.Sprintf("category-%d", category.Id)
		if req.Form.Has(cc) && req.Form.Get(cc) == "on" {
			PostCategories = append(PostCategories, category)
		}
	}
	// Validate user is logged in
	if !user.LoggedIn {
		data.Error = *(&Error{}).Consume(ErrorPostPermissionDenied)
		data.SetView("user_login_view").WriteResponse(res)
		return
	}
	// Create post object
	post := Post{
		Title:      title,
		Body:       body,
		UserId:     user.Id,
		Categories: PostCategories,
	}
	// Save post to database
	postId, err := addPost(post)
	if err != nil {
		data.Error = *(&Error{}).Consume(err)
		data.SetView("post_create_view").WriteResponse(res)
		return
	}
	redirectURL := fmt.Sprintf("/post?id=%d", postId)
	http.Redirect(res, req, redirectURL, http.StatusSeeOther)
}

func handlePostReaction(res http.ResponseWriter, req *http.Request, user User) {
	err := req.ParseForm()
	if err != nil {
		(&Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	postIdStr := req.URL.Query().Get("id")
	if len(postIdStr) == 0 {
		(&Error{}).Consume(ErrorPostEmptyId).LogAndRespondError(res, user)
		return
	}
	ok, err := regexp.MatchString(`^\d+$`, postIdStr)
	if !ok {
		(&Error{}).Consume(ErrorInvalidPostId).LogAndRespondError(res, user)
		return
	}
	postId, err := strconv.ParseInt(postIdStr, 10, 64)
	if err != nil {
		(&Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	if !user.LoggedIn {
		(&Error{}).Consume(ErrorPostPermissionDenied).LogAndRespondError(res, user)
		return
	}
	if req.FormValue("action") == "like" {
		err = DoLike(user.Id, postId)
		if err != nil {
			(&Error{}).Consume(err).LogAndRespondError(res, user)
			return
		}
	}
	if req.FormValue("action") == "dislike" {
		err = DoDislikePost(user.Id, postId)
		if err != nil {
			(&Error{}).Consume(err).LogAndRespondError(res, user)
			return
		}
	}
	http.Redirect(res, req, "/post?id="+postIdStr, http.StatusSeeOther)
}
