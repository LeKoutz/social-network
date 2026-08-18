package models

import (
	"errors"
	"forum/src/db"
	"forum/src/utils"
)

func (user *UserType) LikeComment(commentId int64) error {
	alreadyLiked, err := HasUserLikedComment(user.Id, commentId)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	if alreadyLiked {
		return db.RemoveLikeFromComment(user.Id, commentId)
	}
	existingDislikeId, err := db.CheckIfUserDislikedComment(user.Id, commentId)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	if existingDislikeId != 0 {
		if err = db.RemoveReaction(existingDislikeId); err != nil {
			return errors.Join(utils.GetFunctionName(), err)
		}
	}
	return db.AddLikeToComment(user.Id, commentId)
}

func (user *UserType) LikePost(postId int64) error {
	alreadyLiked, err := HasUserLikedPost(user.Id, postId)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	if alreadyLiked {
		return db.RemoveLikeFromPost(user.Id, postId)
	}
	existingDislikeId, err := db.CheckIfUserDislikedPost(user.Id, postId)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	if existingDislikeId != 0 {
		if err = db.RemoveReaction(existingDislikeId); err != nil {
			return errors.Join(utils.GetFunctionName(), err)
		}
	}
	return db.AddLikeToPost(user.Id, postId)
}
