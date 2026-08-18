package db

import (
	"database/sql"
	"errors"
	"forum/src/utils"
)

func CheckIfUserLikedPost(userId, postId int64) (int64, error) {
	var existingReactionId int64
	err := db.QueryRow(`
		SELECT id FROM reactions
		WHERE user_id = ? AND post_id = ? AND value = 1
		`, userId, postId).Scan(&existingReactionId)
	if err != nil && err != sql.ErrNoRows {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return 0, err
	}
	return existingReactionId, nil
}

func AddLikeToPost(userId, postId int64) error {
	_, err := db.Exec(`
		INSERT INTO reactions (user_id, post_id, value, timestamp)
		VALUES (?, ?, 1, ?)
		`, userId, postId, utils.GetCurrentTimestamp())
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}

func RemoveLikeFromPost(userId, postId int64) error {
	_, err := db.Exec(`
		DELETE FROM reactions
		WHERE user_id = ? AND post_id = ? AND value = 1
		`, userId, postId)
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	return nil
}

func CheckIfUserLikedComment(userId, commentId int64) (int64, error) {
	var existingReactionId int64
	err := db.QueryRow(`
		SELECT id FROM reactions
		WHERE user_id = ? AND comment_id = ? AND value = 1
		`, userId, commentId).Scan(&existingReactionId)
	if err != nil && err != sql.ErrNoRows {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return 0, err
	}
	return existingReactionId, nil
}

func AddLikeToComment(userId, commentId int64) error {
	_, err := db.Exec(`
		INSERT INTO reactions (user_id, comment_id, value, timestamp)
		VALUES (?, ?, 1, ?)
		`, userId, commentId, utils.GetCurrentTimestamp())
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
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
