package models

import (
	"errors"
	"forum/src/db"
	"forum/src/utils"
)

func (user *UserType) DislikePost(postId int64) error {
	dislikeId, err := db.CheckIfUserDislikedPost(user.Id, postId)
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	if dislikeId != 0 {
		return db.RemoveDislikeFromPost(dislikeId)
	}
	existingLikeId, err := db.CheckIfUserLikedPost(user.Id, postId)
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	if existingLikeId != 0 {
		err = db.RemoveLikeFromPost(user.Id, postId)
		if err != nil {
			return errors.Join(utils.GetFunctionName(), err)
		}
	}
	return db.AddDislikeToPost(user.Id, postId)
}

func (user *UserType) DislikeComment(commentId int64) error {
	dislikeId, err := db.CheckIfUserDislikedComment(user.Id, commentId)
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	if dislikeId != 0 {
		return db.RemoveReaction(dislikeId)
	}
	existingLikeId, err := db.CheckIfUserLikedComment(user.Id, commentId)
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	if existingLikeId != 0 {
		err = db.RemoveReaction(existingLikeId)
		if err != nil {
			return errors.Join(utils.GetFunctionName(), err)
		}
	}
	return db.AddDislikeToComment(user.Id, commentId)
}
