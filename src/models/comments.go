package models

import (
	"errors"
	"forum/src/db"
	"forum/src/utils"
)

type CommentsType []CommentType

func (u *UserType) GetCommentsByUserId() (CommentsType, error) {
	var comments CommentsType
	var err error
	rows, err := db.SelectCommentsByUserId(u.Id)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return CommentsType{}, err
	}
	for _, row := range rows {
		var comment CommentType
		comment.CommentRowType = row
		comments = append(comments, comment)
	}
	return comments, nil
}
