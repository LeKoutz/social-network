package controllers

import "forum/src/models"

func DoLikePost(userId, postId int64) error {
	alreadyLiked, err := models.HasUserLikedPost(userId, postId)
	if err != nil {
		return err
	}
	if alreadyLiked {
		return models.RemoveLikeFromPost(userId, postId)
	}
	existingDislikeId, err := models.CheckIfUserDislikedPost(userId, postId)
	if err != nil {
		return err
	}
	if existingDislikeId != 0 {
		return models.RemoveReaction(existingDislikeId)
	}
	return models.AddLikeToPost(userId, postId)
}

func DoLikeComment(userId, commentId int64) error {
	alreadyLiked, err := models.HasUserLikedComment(userId, commentId)
	if err != nil {
		return err
	}
	if alreadyLiked {
		return models.RemoveLikeFromComment(userId, commentId)
	}
	existingDislikeId, err := models.CheckIfUserDislikedComment(userId, commentId)
	if err != nil {
		return err
	}
	if existingDislikeId != 0 {
		return models.RemoveReaction(existingDislikeId)
	}
	return models.AddLikeToComment(userId, commentId)
}
