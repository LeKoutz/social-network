package db

import (
	"errors"
	"forum/src/utils"
)

type ReactionRowType struct {
	Id              int64
	PostId          int64
	UserId          int64
	CommentId       int64
	TimestampString string
}

func DeleteReactionById(reactionId int64) error {
	_, err := db.Exec(`
		DELETE FROM reactions
		WHERE id = ?
		`, reactionId)
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	return nil
}
