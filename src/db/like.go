package db

import (
	"database/sql"
	"errors"
	"forum/src/utils"
)

func SelectUserLikeFromPost(userId, postId int64) (int64, error) {
	var existingReactionId int64
	err := db.QueryRow(`
		SELECT id FROM reactions
		WHERE user_id = ? AND post_id = ? AND value = 1
		`, userId, postId).Scan(&existingReactionId)
	if err != nil && err != sql.ErrNoRows {
		err = errors.Join(utils.GetFunctionName(), err)
		return 0, err
	}
	return existingReactionId, nil
}

func InsertLikeToPost(userId, postId int64) error {
	_, err := db.Exec(`
		INSERT INTO reactions (user_id, post_id, value, timestamp)
		VALUES (?, ?, 1, ?)
		`, userId, postId, utils.GetCurrentTimestamp())
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	return nil
}

func SelectUserLikeFromComment(userId, commentId int64) (int64, error) {
	var existingReactionId int64
	err := db.QueryRow(`
		SELECT id FROM reactions
		WHERE user_id = ? AND comment_id = ? AND value = 1
		`, userId, commentId).Scan(&existingReactionId)
	if err != nil && err != sql.ErrNoRows {
		err = errors.Join(utils.GetFunctionName(), err)
		return 0, err
	}
	return existingReactionId, nil
}

func InsertLikeToComment(userId, commentId int64) error {
	_, err := db.Exec(`
		INSERT INTO reactions (user_id, comment_id, value, timestamp)
		VALUES (?, ?, 1, ?)
		`, userId, commentId, utils.GetCurrentTimestamp())
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	return nil
}

func RemoveLikeFromComment(userId, commentId int64) error {
	_, err := db.Exec(`
		DELETE FROM reactions
		WHERE user_id = ? AND comment_id = ? AND value = 1
		`, userId, commentId)
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	return nil
}
