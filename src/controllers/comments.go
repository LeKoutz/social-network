package controllers

import (
	"fmt"
	"forum/src/models"
	"net/http"
	"regexp"
	"strconv"
)

func createComment(data models.ResponseStruct) {
	// Parse form data
	err := data.Request.ParseForm()
	if err != nil {
		data.Error = *(&models.Error{}).Consume(err)
		data.SetView("user_register_view").WriteResponse()
		return
	}
	// Get form values
	body := data.Request.FormValue("comment")
	postIdStr := data.Request.FormValue("post_id")
	ok, err := regexp.MatchString(`^\d+$`, postIdStr)
	if !ok {
		(&models.Error{}).Consume(models.ErrorInvalidPostId).LogAndRespondError(data.Response, data.User)
		return
	}
	post_id, err := strconv.ParseInt(postIdStr, 10, 64)
	if err != nil {
		data.Error = *(&models.Error{}).Consume(err)
		data.SetView("post_view").WriteResponse()
		return
	}
	// Validate user is logged in
	if !data.User.LoggedIn {
		data.Error = *(&models.Error{}).Consume(models.ErrorCommentPermissionDenied)
		data.SetView("user_login_view").WriteResponse()
		return
	}
	// Create post object
	comment := models.Comment{
		Body:   body,
		UserId: data.User.Id,
		PostId: post_id,
	}
	// Save post to database
	commentId, err := models.AddComment(comment)
	if err != nil {
		data.Error = *(&models.Error{}).Consume(err)
		data.SetView("post_view").WriteResponse()
		return
	}
	redirectURL := fmt.Sprintf("/post?id=%d#%d", post_id, commentId)
	// Redirect to the post's page
	http.Redirect(data.Response, data.Request, redirectURL, http.StatusSeeOther)
}

func handleCommentReaction(data models.ResponseStruct) {
	err := data.Request.ParseForm()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	commentIdStr := data.Request.URL.Query().Get("id")
	if len(commentIdStr) == 0 {
		(&models.Error{}).Consume(models.ErrorCommentEmptyId).LogAndRespondError(data.Response, data.User)
		return
	}
	ok, err := regexp.MatchString(`^\d+$`, commentIdStr)
	if !ok {
		(&models.Error{}).Consume(models.ErrorInvalidCommentId).LogAndRespondError(data.Response, data.User)
		return
	}
	commentId, err := strconv.ParseInt(commentIdStr, 10, 64)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	if !data.User.LoggedIn {
		(&models.Error{}).Consume(models.ErrorCommentPermissionDenied).LogAndRespondError(data.Response, data.User)
		return
	}
	if data.Request.FormValue("action") == "like" {
		err = DoLikeComment(data.User.Id, commentId)
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
			return
		}
	}
	if data.Request.FormValue("action") == "dislike" {
		err = DoDislikeComment(data.User.Id, commentId)
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
			return
		}
	}
	postId := data.Request.FormValue("post-id")
	http.Redirect(data.Response, data.Request, "/post?id="+postId+"#"+commentIdStr, http.StatusSeeOther)
}
