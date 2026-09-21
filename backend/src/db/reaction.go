package db

import "forum/src/ferror"

type ReactionRowType struct {
	Id              int64
	PostId          int64
	UserId          int64
	CommentId       int64
	Timestamp       string
}

func DeleteReactionById(reactionId int64) error {
	_, err := db.Exec(`
		DELETE FROM reactions
		WHERE id = ?
		`, reactionId)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}
