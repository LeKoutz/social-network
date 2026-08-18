package models

import (
	"errors"
	"forum/src/db"
	"forum/src/utils"
)

func (user *UserType) DislikePost(postId int64) error {
	dislikeId, err := db.SelectUserDislikeFromPost(user.Id, postId)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	if dislikeId != 0 {
		return db.DeleteReactionById(dislikeId)
	}
	existingLikeId, err := db.SelectUserLikeFromPost(user.Id, postId)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	if existingLikeId != 0 {
		err = db.DeleteReactionById(existingLikeId)
		if err != nil {
			if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return err
		}
	}
	return db.InsertDislikeToPost(user.Id, postId)
}

func (user *UserType) DislikeComment(commentId int64) error {
	dislikeId, err := db.SelectUserDislikeFromComment(user.Id, commentId)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	if dislikeId != 0 {
		return db.DeleteReactionById(dislikeId)
	}
	existingLikeId, err := db.SelectUserLikeFromComment(user.Id, commentId)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	if existingLikeId != 0 {
		err = db.DeleteReactionById(existingLikeId)
		if err != nil {
			if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return err
		}
	}
	return db.InsertDislikeToComment(user.Id, commentId)
}
