package db

import (
	"database/sql"
	"errors"
	"forum/src/utils"
)

func CheckIfUserDislikedPost(userId, postId int64) (int64, error) {
	var existingDislikeId int64
	err := db.QueryRow(`
		SELECT id FROM reactions
		WHERE user_id = ? AND post_id = ? AND value = 2
		`, userId, postId).Scan(&existingDislikeId)
	if err != nil && err != sql.ErrNoRows {
		err = errors.Join(utils.GetFunctionName(), err)
		return 0, err
	}
	return existingDislikeId, nil
}

func AddDislikeToPost(userId, postId int64) error {
	_, err := db.Exec(`
		INSERT INTO reactions (user_id, post_id, value, timestamp)
		VALUES (?, ?, 2, ?)
		`, userId, postId, utils.GetCurrentTimestamp())
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	return nil
}

func RemoveDislikeFromPost(dislikeId int64) error {
	_, err := db.Exec(`
		DELETE FROM reactions
		WHERE id = ?
		`, dislikeId)
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	return nil
}

func CheckIfUserDislikedComment(userId, commentId int64) (int64, error) {
	var existingDislikeId int64
	err := db.QueryRow(`
		SELECT id FROM reactions
		WHERE user_id = ? AND comment_id = ? AND value = 2
		`, userId, commentId).Scan(&existingDislikeId)
	if err != nil && err != sql.ErrNoRows {
		err = errors.Join(utils.GetFunctionName(), err)
		return 0, err
	}
	return existingDislikeId, nil
}

func AddDislikeToComment(userId, commentId int64) error {
	_, err := db.Exec(`
		INSERT INTO reactions (user_id, comment_id, value, timestamp)
		VALUES (?, ?, 2, ?)
		`, userId, commentId, utils.GetCurrentTimestamp())
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	return nil
}
