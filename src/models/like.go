package models

import (
	"database/sql"
	"errors"
	"forum/src/utils"
)

type Like struct {
	PostId    int64
	UserId    int64
	CommentId int64
}

func CheckIfUserLikedPost(userId, postId int64) (int64, error) {
	var existingReactionId int64
	err := DB.QueryRow(`
		SELECT id FROM reactions
		WHERE user_id = ? AND post_id = ? AND value = 1
		`, userId, postId).Scan(&existingReactionId)
	if err != nil && err != sql.ErrNoRows {
		err = errors.Join(utils.GetFunctionName(), err)
		return 0, err
	}
	return existingReactionId, nil
}

func CheckIfUserDislikedPost(userId, postId int64) (int64, error) {
	var existingDislikeId int64
	err := DB.QueryRow(`
		SELECT id FROM reactions
		WHERE user_id = ? AND post_id = ? AND value = 2
		`, userId, postId).Scan(&existingDislikeId)
	if err != nil && err != sql.ErrNoRows {
		err = errors.Join(utils.GetFunctionName(), err)
		return 0, err
	}
	return existingDislikeId, nil
}

func AddLikeToPost(userId, postId int64) error {
	_, err := DB.Exec(`
		INSERT INTO reactions (user_id, post_id, value)
		VALUES (?, ?, 1)
		`, userId, postId)
	if err != nil {
		return err
	}
	return nil
}

func RemoveDislikeFromPost(dislikeId int64) error {
	_, err := DB.Exec(`
		DELETE FROM reactions
		WHERE id = ?
		`, dislikeId)
	if err != nil {
		return err
	}
	return nil
}

func AddDislikeToPost(userId, postId int64) error {
	_, err := DB.Exec(`
		INSERT INTO reactions (user_id, post_id, value)
		VALUES (?, ?, 2)
		`, userId, postId)
	if err != nil {
		return err
	}
	return nil
}

func RemoveLikeFromPost(userId, postId int64) error {
	_, err := DB.Exec(`
		DELETE FROM reactions
		WHERE user_id = ? AND post_id = ? AND value = 1
		`, userId, postId)
	if err != nil {
		return err
	}
	return nil
}

func CheckIfUserLikedComment(userId, commentId int64) (int64, error) {
	var existingReactionId int64
	err := DB.QueryRow(`
		SELECT id FROM reactions
		WHERE user_id = ? AND comment_id = ? AND value = 1
		`, userId, commentId).Scan(&existingReactionId)
	if err != nil && err != sql.ErrNoRows {
		err = errors.Join(utils.GetFunctionName(), err)
		return 0, err
	}
	return existingReactionId, nil
}

func CheckIfUserDislikedComment(userId, commentId int64) (int64, error) {
	var existingDislikeId int64
	err := DB.QueryRow(`
		SELECT id FROM reactions
		WHERE user_id = ? AND comment_id = ? AND value = 2
		`, userId, commentId).Scan(&existingDislikeId)
	if err != nil && err != sql.ErrNoRows {
		err = errors.Join(utils.GetFunctionName(), err)
		return 0, err
	}
	return existingDislikeId, nil
}

func AddLikeToComment(userId, commentId int64) error {
	_, err := DB.Exec(`
		INSERT INTO reactions (user_id, comment_id, value)
		VALUES (?, ?, 1)
		`, userId, commentId)
	return err
}

func AddDislikeToComment(userId, commentId int64) error {
	_, err := DB.Exec(`
		INSERT INTO reactions (user_id, comment_id, value)
		VALUES (?, ?, 2)
		`, userId, commentId)
	return err
}

func RemoveLikeFromComment(userId, commentId int64) error {
	_, err := DB.Exec(`
		DELETE FROM reactions
		WHERE user_id = ? AND comment_id = ? AND value = 1
		`, userId, commentId)
	return err
}
