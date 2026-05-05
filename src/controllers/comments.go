package controllers

import (
	"fmt"
	"forum/src/models"
	"forum/src/utils"
	"forum/src/views"
	"net/http"
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
	post_id, err := utils.StringToInt64(postIdStr)
	if err != nil {
		err = models.ErrorInvalidPostId
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
	createComment(data)
}

func handleCommentReaction(data models.ResponseStruct) {
	var err error
	commentIdStr := data.Request.FormValue("comment-id")
	if len(commentIdStr) == 0 {
		(&models.Error{}).Consume(models.ErrorCommentEmptyId).LogAndRespondError(data.Response, data.User)
		return
	}
	commentId, err := utils.StringToInt64(commentIdStr)
	if err != nil {
		err = models.ErrorInvalidCommentId
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
	notification := models.Notification{
		UserId:    comment.UserId,
		ActorId:   data.User.Id,
		CommentId: comment.Id,
		PostId:    comment.PostId,
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
	redirectURL := fmt.Sprintf("/post/view/%d#comment-%d", comment.PostId, comment.Id)
	http.Redirect(data.Response, data.Request, redirectURL, http.StatusSeeOther)
}

func handleCommentDelete(data models.ResponseStruct) {
	if !data.User.LoggedIn {
		(&models.Error{}).Consume(models.ErrorCommentPermissionDenied).LogAndRespondError(data.Response, data.User)
		return
	}
	if data.Request.Method != http.MethodPost {
		(&models.Error{}).Consume(models.ErrorMethodNotAllowed).LogAndRespondError(data.Response, data.User)
		return
	}
	commentIdStr := data.Request.FormValue("comment-id")
	if len(commentIdStr) == 0 {
		(&models.Error{}).Consume(models.ErrorCommentEmptyId).LogAndRespondError(data.Response, data.User)
		return
	}
	commentId, err := utils.StringToInt64(commentIdStr)
	if err != nil {
		err = models.ErrorInvalidCommentId
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	comment := models.Comment{Id: commentId}
	err = comment.GetCommentById()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	if comment.UserId != data.User.Id {
		(&models.Error{}).Consume(models.ErrorCommentPermissionDenied).LogAndRespondError(data.Response, data.User)
		return
	}
	err = comment.Delete()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	redirectURL := fmt.Sprintf("/post/view/%d", comment.PostId)
	http.Redirect(data.Response, data.Request, redirectURL, http.StatusSeeOther)
}

func handleCommentEdit(data models.ResponseStruct) {
	var err error
	if !data.User.LoggedIn {
		(&models.Error{}).Consume(models.ErrorCommentPermissionDenied).LogAndRespondError(data.Response, data.User)
		return
	}
	if data.Request.Method != http.MethodPost {
		(&models.Error{}).Consume(models.ErrorMethodNotAllowed).LogAndRespondError(data.Response, data.User)
		return
	}
	err = validateFormCommentEdit(&data)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	err = verifyCommentPostAssociation(&data)
	if err != nil {
		(&models.Error{}).Consume(models.ErrorMethodNotAllowed).LogAndRespondError(data.Response, data.User)
		return
	}
	if data.Request.FormValue("save-comment") == "1" {
		err = updateCommentFromForm(&data)
		if err != nil {
			(&models.Error{}).Consume(models.ErrorMethodNotAllowed).LogAndRespondError(data.Response, data.User)
			return
		}
	}
	err = showEditComment(&data)
	if err != nil {
		(&models.Error{}).Consume(models.ErrorMethodNotAllowed).LogAndRespondError(data.Response, data.User)
		return
	}
	views.PostView(data)
}

func validateFormCommentEdit(data *models.ResponseStruct) error {
	var err error
	var post models.Post
	var comment models.Comment
	commentIdStr := data.Request.FormValue("comment-id")
	if len(commentIdStr) == 0 {
		return models.ErrorCommentEmptyId
	}
	commentId, err := utils.StringToInt64(commentIdStr)
	if err != nil {
		return models.ErrorInvalidCommentId
	}
	postIdStr := data.Request.FormValue("post-id")
	if len(postIdStr) == 0 {
		return models.ErrorPostEmptyId
	}
	postId, err := utils.StringToInt64(postIdStr)
	if err != nil {
		return models.ErrorInvalidPostId
	}
	comment = models.Comment{Id: commentId}
	post.Comments = models.Comments{comment}
	post.Id = postId
	data.Posts = models.Posts{post}
	return nil
}

func verifyCommentPostAssociation(data *models.ResponseStruct) error {
	var err error
	post := &data.Posts[0]
	comment := &data.Posts[0].Comments[0]
	err = comment.GetCommentById()
	if err != nil {
		return err
	}
	if comment.PostId != post.Id {
		return models.ErrorInvalidPostId
	}
	// Check your priviledge
	if comment.UserId != data.User.Id {
		return models.ErrorCommentPermissionDenied
	}
	return nil
}

func showEditComment(data *models.ResponseStruct) error {
	var err error
	comment := data.Posts[0].Comments[0]
	err = getPostDataById(data)
	if err != nil {
		return err
	}
	data.EditCommentId = comment.Id
	return nil
}

func updateCommentFromForm(data *models.ResponseStruct) error {
	var err error
	comment := &data.Posts[0].Comments[0]
	comment.Body = data.Request.FormValue("comment")
	err = comment.Update()
	if err != nil {
		return err
	}
	redirectURL := fmt.Sprintf("/post/view/%d#comment-%d", comment.PostId, comment.Id)
	http.Redirect(data.Response, data.Request, redirectURL, http.StatusSeeOther)
	return nil
}
