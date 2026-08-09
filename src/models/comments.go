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
	comment_rows, err := db.GetCommentsByUserId(u.Id)
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return CommentsType{}, err
	}
	for _, comment_row := range comment_rows {
		var comment CommentType
		comment.FromCommentRowType(comment_row)
		t, err := utils.ConvertStringToTime(comment.TimestampString)
		if err != nil {
			if config.Debug {
				err = errors.Join(utils.GetFunctionName(), err)
			}
			return CommentsType{}, err
		}
		comment.TimestampString = utils.ConvertTimeToString(t)
		comments = append(comments, comment)
	}
	return comments, nil
}
