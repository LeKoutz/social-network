package controllers

import (
	"fmt"
	"forum/src/models"
	"forum/src/utils"
	"forum/src/views"
	"net/http"
	"regexp"
	"strconv"
)

func createComment(data models.ResponseStruct) {
	// Validate user is logged in
	if !data.User.LoggedIn {
		data.Error.Consume(models.ErrorCommentPermissionDenied)
		views.UserLogin(data)
		return
	}
	// Parse form data
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
		data.Error.Consume(err)
		views.ErrorView(data)
		return
	}
	// Create post object
	comment := models.Comment{
		Body:   body,
		UserId: data.User.Id,
		PostId: post_id,
	}
	// Save post to database
	commentId, err := comment.Add()
	if err != nil {
		data.Error.Consume(err)
		views.ErrorView(data)
		return
	}
	// Create notification for the post's author
	post := models.Post{Id: post_id}
	err = post.GetById()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	notification := models.Notification{
		UserId:    post.User.Id,
		ActorId:   data.User.Id,
		Type:      "comment",
		PostId: 	post_id,
		CommentId:  commentId,
		TimestampString: utils.GetCurrentTimestamp(),
	}
	if data.User.Id != post.User.Id {
		err = notification.Add()
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
			return
		}
	}
	redirectURL := fmt.Sprintf("/post/view/%d#comment-%d", post_id, commentId)
	// Redirect to the post's page
	http.Redirect(data.Response, data.Request, redirectURL, http.StatusSeeOther)
}

func handleCommentCreate(data models.ResponseStruct) {
	if data.Request.Method != http.MethodPost {
		(&models.Error{}).Consume(models.ErrorMethodNotAllowed).LogAndRespondError(data.Response, data.User)
		return
	}
	err := data.Request.ParseForm()
	if err != nil {
		data.Error.Consume(err)
		data.SetView("error_view").WriteResponse()
		return
	}
	createComment(data)
}

func handleCommentReaction(data models.ResponseStruct) {
	err := data.Request.ParseForm()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	commentIdStr := data.Request.FormValue("comment-id")
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
	comment := models.Comment{Id: commentId}
	err = comment.GetCommentById()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	postIdStr := data.Request.FormValue("post-id")
	if len(postIdStr) == 0 {
		(&models.Error{}).Consume(models.ErrorPostEmptyId).LogAndRespondError(data.Response, data.User)
		return
	}
	ok, err = regexp.MatchString(`^\d+$`, postIdStr)
	if !ok {
		(&models.Error{}).Consume(models.ErrorInvalidPostId).LogAndRespondError(data.Response, data.User)
		return
	}
	postId, err := strconv.ParseInt(postIdStr, 10, 64)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	notification := models.Notification{
		UserId:    comment.UserId,
		ActorId:   data.User.Id,
		CommentId: comment.Id,
		PostId:    postId,
	}
	if data.Request.FormValue("action") == "like" {
		err = DoLikeComment(data.User.Id, commentId)
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
			return
		}
		if data.User.Id != comment.UserId {
			notification.Type = "commentLike"
			notification.TimestampString = utils.GetCurrentTimestamp()
			err = notification.Add()
			if err != nil {
				(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
				return
			}
		}
	}
	if data.Request.FormValue("action") == "dislike" {
		err = DoDislikeComment(data.User.Id, commentId)
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
			return
		}
		if data.User.Id != comment.UserId {
			notification.Type = "commentDislike"
			notification.TimestampString = utils.GetCurrentTimestamp()
			err = notification.Add()
			if err != nil {
				(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
				return
			}
		}
	}
	http.Redirect(data.Response, data.Request, "/post/view/"+postIdStr+"#comment-"+commentIdStr, http.StatusSeeOther)
}
