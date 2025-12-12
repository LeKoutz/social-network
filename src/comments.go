package forum

import (
	"fmt"
	"net/http"
	"strconv"
)

type Comment struct {
	Id              int64
	PostId          int64
	UserId          int64
	Body            string
	Timestamp       int64
	TimestampString string
	Likes           int
	Liked           bool
	Dislikes        int
	Disliked        bool
	Username        string
}

type Comments []Comment

func (c *Comment) validateComment() error {
	if len(c.Body) == 0 {
		return ErrorCommentEmpty
	}
	if len(c.Body) > 1000 {
		return ErrorCommentTooLong
	}
	return nil
}

func ReturnMockComments() Comments {
	return Comments{
		{
			Id:       1,
			Body:     "mpla mpla",
			Likes:    2,
			Disliked: true,
			Dislikes: 1,
		},
	}
}

func createComment(res http.ResponseWriter, req *http.Request, user User) {
	data := ResponseStruct{}
	data.Init().SetUser(user)

	// Parse form data
	err := req.ParseForm()
	if err != nil {
		data.Error = *(&Error{}).Consume(err)
		data.SetView("user_register_view").WriteResponse(res)
		return
	}

	// Get form values
	body := req.FormValue("comment")
	post_id, err := strconv.ParseInt(req.FormValue("post_id"), 10, 64)
	if err != nil {
		data.Error = *(&Error{}).Consume(err)
		data.SetView("post_view").WriteResponse(res)
		return
	}

	// Validate user is logged in
	if !user.LoggedIn {
		data.Error = *(&Error{}).Consume(ErrorCommentPermissionDenied)
		data.SetView("user_login_view").WriteResponse(res)
		return
	}

	// Create post object
	comment := Comment{
		Body:   body,
		UserId: user.Id,
		PostId: post_id,
	}

	// Save post to database
	commentId, err := addComment(comment)
	if err != nil {
		data.Error = *(&Error{}).Consume(err)
		data.SetView("post_view").WriteResponse(res)
		return
	}

	redirectURL := fmt.Sprintf("/post?id=%d#%d", post_id, commentId)
	// Redirect to the post's page
	http.Redirect(res, req, redirectURL, http.StatusSeeOther)
}

func handleCommentReaction(res http.ResponseWriter, req *http.Request, user User) {
	err := req.ParseForm()
	if err != nil {
		(&Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	commentIdStr := req.URL.Query().Get("id")
	if len(commentIdStr) == 0 {
		(&Error{}).Consume(ErrorCommentEmptyId).LogAndRespondError(res, user)
		return
	}
	commentId, err := strconv.ParseInt(commentIdStr, 10, 64)
	if err != nil {
		(&Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	if !user.LoggedIn {
		(&Error{}).Consume(ErrorCommentPermissionDenied).LogAndRespondError(res, user)
		return
	}
	if req.FormValue("like") == "on" {
		err = DoLikeComment(user.Id, commentId)
		if err != nil {
			(&Error{}).Consume(err).LogAndRespondError(res, user)
			return
		}
	}
	if req.FormValue("dislike") == "on" {
		err = DoDislikeComment(user.Id, commentId)
		if err != nil {
			(&Error{}).Consume(err).LogAndRespondError(res, user)
			return
		}
	}
	http.Redirect(res, req, "/", http.StatusSeeOther)
}
