package forum

import (
	"fmt"
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
	if p.Categories.IsEmpty() {
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
	post.Likes, err = getLikesCountByPostId(id_int)
	if err != nil {
		(&Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	post.Liked, err = hasUserAlreadyLikedPost(user.Id, id_int)
	if err != nil {
		(&Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	post.Dislikes, err = getDislikesCountByPostId(id_int)
	if err != nil {
		(&Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	post.Disliked, err = hasUserAlreadyDislikedPost(user.Id, id_int)
	if err != nil {
		(&Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
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
	for i := range posts {
		posts[i].Likes, err = getLikesCountByPostId(posts[i].Id)
		if err != nil {
			(&Error{}).Consume(err).LogAndRespondError(res, user)
			return
		}
		posts[i].Liked, err = hasUserAlreadyLikedPost(user.Id, posts[i].Id)
		if err != nil {
			(&Error{}).Consume(err).LogAndRespondError(res, user)
			return
		}
		posts[i].Dislikes, err = getDislikesCountByPostId(posts[i].Id)
		if err != nil {
			(&Error{}).Consume(err).LogAndRespondError(res, user)
			return
		}
		posts[i].Disliked, err = hasUserAlreadyDislikedPost(user.Id, posts[i].Id)
		if err != nil {
			(&Error{}).Consume(err).LogAndRespondError(res, user)
			return
		}
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

	err := req.ParseForm()
	if err != nil {
		data.Error = *(&Error{}).Consume(err)
		respondView(res, "user_register_view", data)
		return
	}
	// Get form values
	title := req.FormValue("title")
	body := req.FormValue("body")
	categories, err := getAllCategories()
	if err != nil {
		data.Error = *(&Error{}).Consume(err)
		respondView(res, "post_create_view", data)
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
		respondView(res, "user_login_view", data)
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
		respondView(res, "post_create_view", data)
		return
	}
	postIdStr := strconv.Itoa(postId)
	redirectURL := "/post?id=" + postIdStr
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
	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		(&Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	if !user.LoggedIn {
		(&Error{}).Consume(ErrorPostPermissionDenied).LogAndRespondError(res, user)
		return
	}
	if req.FormValue("like") == "on" {
		err = DoLike(user.Id, postId)
		if err != nil {
			(&Error{}).Consume(err).LogAndRespondError(res, user)
			return
		}
	}
	if req.FormValue("like") == "" {
		err = UndoLike(user.Id, postId)
		if err != nil {
			(&Error{}).Consume(err).LogAndRespondError(res, user)
			return
		}
	}
	if req.FormValue("dislike") == "on" {
		err = DoDislikePost(user.Id, postId)
		if err != nil {
			(&Error{}).Consume(err).LogAndRespondError(res, user)
			return
		}
	}
	if req.FormValue("dislike") == "" {
		err = UndoDislike(user.Id, postId)
		if err != nil {
			(&Error{}).Consume(err).LogAndRespondError(res, user)
			return
		}
	}
	http.Redirect(res, req, "/post?id="+postIdStr, http.StatusSeeOther)
}
