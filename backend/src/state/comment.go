package state

import "forum/src/models"

func (r *State) InitComment() {
	var comment models.CommentType
	if r.Posts == nil {
		r.InitPost()
	}
	r.EditPost().Comments = models.CommentsType{comment}
}

func (r *State) SetComment(comment models.CommentType) {
	r.EditPost().Comments = models.CommentsType{comment}
}

func (r *State) GetComment() models.CommentType {
	return *r.EditComment()
}

func (r *State) SetEditCommentId(id int64) {
	r.EditCommentId = id
}

func (r *State) EditComment() *models.CommentType {
	if r.Posts == nil {
		r.InitPost()
	}
	if r.EditPost().Comments == nil {
		r.InitComment()
	}
	return &r.Posts[0].Comments[0]
}
