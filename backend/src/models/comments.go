package models

import (
	"forum/src/db"
	"forum/src/ferror"
)

type CommentsType []CommentType

func (u *UserType) GetCommentsByUserId() (CommentsType, error) {
	var comments CommentsType
	var err error
	rows, err := db.SelectCommentsByUserId(u.Id)
	if err != nil {
		return CommentsType{}, ferror.ReturnErr(err)
	}
	for _, row := range rows {
		var comment CommentType
		comment.CommentRowType = row
		comments = append(comments, comment)
	}
	return comments, nil
}
