package controllers

import "forum/src/models"

type Dislikes []Dislike

type Dislike struct {
	PostId    int64
	UserId    int64
	CommentId int64
}

func DoDislikePost(userId, postId int64) error {
	dislikeId, err := models.CheckIfUserDislikedPost(userId, postId)
	if err != nil {
		return err
	}
	if dislikeId != 0 {
		return models.RemoveDislikeFromPost(dislikeId)
	}
	existingLikeId, err := models.CheckIfUserLikedPost(userId, postId)
	if err != nil {
		return err
	}
	if existingLikeId != 0 {
		err = models.RemoveLikeFromPost(userId, postId)
		if err != nil {
			return err
		}
	}
	return models.AddDislikeToPost(userId, postId)
}

func DoDislikeComment(userId, commentId int64) error {
	dislikeId, err := models.CheckIfUserDislikedComment(userId, commentId)
	if err != nil {
		return err
	}
	if dislikeId != 0 {
		return models.RemoveReaction(dislikeId)
	}
	existingLikeId, err := models.CheckIfUserLikedComment(userId, commentId)
	if err != nil {
		return err
	}
	if existingLikeId != 0 {
		err = models.RemoveReaction(existingLikeId)
		if err != nil {
			return err
		}
	}
	return models.AddDislikeToComment(userId, commentId)
}
