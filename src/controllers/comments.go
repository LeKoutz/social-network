package controllers

import (
	"fmt"
	"forum/src/models"
	"net/http"
	"regexp"
	"strconv"
)

func createComment(res http.ResponseWriter, req *http.Request, user models.User) {
	data := models.ResponseStruct{}
	data.Init().SetUser(user).SetResponse(res)
	// Parse form data
	err := req.ParseForm()
	if err != nil {
		data.Error = *(&models.Error{}).Consume(err)
		data.SetView("user_register_view").WriteResponse()
		return
	}
	// Get form values
	body := req.FormValue("comment")
	postIdStr := req.FormValue("post_id")
	ok, err := regexp.MatchString(`^\d+$`, postIdStr)
	if !ok {
		(&models.Error{}).Consume(models.ErrorInvalidPostId).LogAndRespondError(res, user)
		return
	}
	post_id, err := strconv.ParseInt(postIdStr, 10, 64)
	if err != nil {
		data.Error = *(&models.Error{}).Consume(err)
		data.SetView("post_view").WriteResponse()
		return
	}
	// Validate user is logged in
	if !user.LoggedIn {
		data.Error = *(&models.Error{}).Consume(models.ErrorCommentPermissionDenied)
		data.SetView("user_login_view").WriteResponse()
		return
	}
	// Create post object
	comment := models.Comment{
		Body:   body,
		UserId: user.Id,
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
	http.Redirect(res, req, redirectURL, http.StatusSeeOther)
}

func handleCommentReaction(res http.ResponseWriter, req *http.Request, user models.User) {
	err := req.ParseForm()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	commentIdStr := req.URL.Query().Get("id")
	if len(commentIdStr) == 0 {
		(&models.Error{}).Consume(models.ErrorCommentEmptyId).LogAndRespondError(res, user)
		return
	}
	ok, err := regexp.MatchString(`^\d+$`, commentIdStr)
	if !ok {
		(&models.Error{}).Consume(models.ErrorInvalidCommentId).LogAndRespondError(res, user)
		return
	}
	commentId, err := strconv.ParseInt(commentIdStr, 10, 64)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	if !user.LoggedIn {
		(&models.Error{}).Consume(models.ErrorCommentPermissionDenied).LogAndRespondError(res, user)
		return
	}
	if req.FormValue("action") == "like" {
		err = DoLikeComment(user.Id, commentId)
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(res, user)
			return
		}
	}
	if req.FormValue("action") == "dislike" {
		err = DoDislikeComment(user.Id, commentId)
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(res, user)
			return
		}
	}
	postId := req.FormValue("post-id")
	http.Redirect(res, req, "/post?id="+postId+"#"+commentIdStr, http.StatusSeeOther)
}
