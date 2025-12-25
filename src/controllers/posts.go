package controllers

import (
	"fmt"
	"forum/src/models"
	"net/http"
	"regexp"
	"strconv"
)

func showPost(res http.ResponseWriter, req *http.Request, user models.User) {
	data := models.ResponseStruct{}
	data.Init().SetUser(user).SetResponse(res)
	id := req.URL.Query().Get("id")
	if len(id) == 0 {
		(&models.Error{}).Consume(models.ErrorPostEmptyId).LogAndRespondError(res, user)
		return
	}
	ok, err := regexp.MatchString(`^\d+$`, id)
	if !ok {
		(&models.Error{}).Consume(models.ErrorInvalidPostId).LogAndRespondError(res, user)
		return
	}
	id_int, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	var post models.Post
	post.Id = id_int
	err = post.GetById()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	comments, err := models.GetCommentsByPostId(id_int)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	for i := range comments {
		comments[i].GetReactions()
		comments[i].GetReactionsByUserId(user.Id)
	}
	post.Comments = comments
	categories, err := models.GetCategoriesByPostId(post.Id)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	post.Categories = categories
	err = post.GetReactions()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	err = post.GetReactionsByUserId(user.Id)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	data.Posts = models.Posts{post}
	data.SetView("post_view").WriteResponse()
}

func showPosts(res http.ResponseWriter, _ *http.Request, user models.User) {
	data := models.ResponseStruct{}
	data.Init().SetResponse(res)
	posts, err := models.GetAllPosts()
	if err != nil {
		data.Error = *(&models.Error{}).Consume(err)
		(&models.Error{}).Consume(err).RespondError(res, user)
		return
	}
	for i := range posts {
		err = posts[i].GetReactions()
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(res, user)
			return
		}
		err = posts[i].GetReactionsByUserId(user.Id)
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(res, user)
			return
		}
	}
	data.SetPosts(posts).SetUser(user).SetView("posts_view").WriteResponse()
}

func createPost(res http.ResponseWriter, req *http.Request, user models.User) {
	data := models.ResponseStruct{}
	data.Init().SetUser(user).SetResponse(res)
	if !user.LoggedIn {
		data.Error = *(&models.Error{}).Consume(models.ErrorPostPermissionDenied)
		data.SetView("user_login_view").WriteResponse()
		return
	}
	err := req.ParseForm()
	if err != nil {
		data.Error = *(&models.Error{}).Consume(err)
		data.SetView("post_create_view").WriteResponse()
		return
	}
	// Get form values
	title := req.FormValue("title")
	body := req.FormValue("body")
	categories, err := models.GetAllCategories()
	if err != nil {
		data.Error = *(&models.Error{}).Consume(err)
		data.SetView("post_create_view").WriteResponse()
		return
	}
	var PostCategories models.Categories
	for _, category := range categories {
		cc := fmt.Sprintf("category-%d", category.Id)
		if req.Form.Has(cc) && req.Form.Get(cc) == "on" {
			PostCategories = append(PostCategories, category)
		}
	}
	post := models.Post{
		Title:      title,
		Body:       body,
		UserId:     user.Id,
		Categories: PostCategories,
	}
	postId, err := post.Add()
	if err != nil {
		data.Error = *(&models.Error{}).Consume(err)
		data.SetView("post_create_view").WriteResponse()
		return
	}
	redirectURL := fmt.Sprintf("/post?id=%d", postId)
	http.Redirect(res, req, redirectURL, http.StatusSeeOther)
}

func handlePostReaction(res http.ResponseWriter, req *http.Request, user models.User) {
	err := req.ParseForm()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	postIdStr := req.URL.Query().Get("id")
	if len(postIdStr) == 0 {
		(&models.Error{}).Consume(models.ErrorPostEmptyId).LogAndRespondError(res, user)
		return
	}
	ok, err := regexp.MatchString(`^\d+$`, postIdStr)
	if !ok {
		(&models.Error{}).Consume(models.ErrorInvalidPostId).LogAndRespondError(res, user)
		return
	}
	postId, err := strconv.ParseInt(postIdStr, 10, 64)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	if !user.LoggedIn {
		(&models.Error{}).Consume(models.ErrorPostPermissionDenied).LogAndRespondError(res, user)
		return
	}
	if req.FormValue("action") == "like" {
		err = DoLikePost(user.Id, postId)
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(res, user)
			return
		}
	}
	if req.FormValue("action") == "dislike" {
		err = DoDislikePost(user.Id, postId)
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(res, user)
			return
		}
	}
	http.Redirect(res, req, "/post?id="+postIdStr, http.StatusSeeOther)
}
