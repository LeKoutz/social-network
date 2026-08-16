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
	rows, err := db.GetCommentsByUserId(u.Id)
	if err != nil {
		if config.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return CommentsType{}, err
	}
	for _, row := range rows {
		var comment CommentType
		t, err := utils.ConvertStringToTime(row.TimestampString)
		if err != nil {
			if config.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return CommentsType{}, err
		}
		row.TimestampString = utils.ConvertTimeToString(t)
		comment.CommentRowType = row
		comments = append(comments, comment)
	}
	return comments, nil
}
