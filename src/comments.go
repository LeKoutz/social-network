package forum

import (
	"fmt"
	"net/http"
	"strconv"
)

type Comment struct {
	Id              int
	PostId          int
	UserId          int
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
	body := req.FormValue("comment")
	post_id, err := strconv.Atoi(req.FormValue("post_id"))
	if err != nil {
		data.Error = *(&Error{}).Consume(err)
		respondView(res, "post_view", data)
		return
	}

	// Validate user is logged in
	if !user.LoggedIn {
		data.Error = *(&Error{}).Consume(ErrorCommentPermissionDenied)
		respondView(res, "user_login_view", data)
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
		respondView(res, "post_view", data)
		return
	}

	commentIdStr := strconv.Itoa(commentId)
	redirectURL := fmt.Sprintf("/post?id=%d#%s", post_id, commentIdStr)
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
	commentId, err := strconv.Atoi(commentIdStr)
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
