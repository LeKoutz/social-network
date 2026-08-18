package db

import (
	"errors"
	"forum/src/utils"
)

type CommentRowsType []CommentRowType

func SelectCommentsByUserId(id int64) (CommentRowsType, error) {
	var comments CommentRowsType
	rows, err := db.Query(`
	SELECT id, post_id, body, timestamp, user_id
	FROM comments
	WHERE user_id = ?`, id)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return comments, err
	}
	for rows.Next() {
		var comment CommentRowType
		var ts string
		err = rows.Scan(&comment.Id, &comment.PostId, &comment.Body, &ts, &comment.UserId)
		if err != nil {
			if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return comments, err
		}
		t, err := utils.ConvertStringToTime(ts)
		if err != nil {
			if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return comments, err
		}
		comment.TimestampString = utils.ConvertTimeToString(t)
		comments = append(comments, comment)
	}
	return comments, nil
}
