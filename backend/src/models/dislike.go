package models

import (
	"forum/src/db"
	"forum/src/ferror"
)

func (user *UserType) DislikePost(postId int64) error {
	dislikeId, err := db.SelectUserDislikeFromPost(user.Id, postId)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	if dislikeId != 0 {
		return db.DeleteReactionById(dislikeId)
	}
	existingLikeId, err := db.SelectUserLikeFromPost(user.Id, postId)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	if existingLikeId != 0 {
		err = db.DeleteReactionById(existingLikeId)
		if err != nil {
			return ferror.ReturnErr(err)
		}
	}
	return db.InsertDislikeToPost(user.Id, postId)
}

func (user *UserType) DislikeComment(commentId int64) error {
	dislikeId, err := db.SelectUserDislikeFromComment(user.Id, commentId)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	if dislikeId != 0 {
		return db.DeleteReactionById(dislikeId)
	}
	existingLikeId, err := db.SelectUserLikeFromComment(user.Id, commentId)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	if existingLikeId != 0 {
		err = db.DeleteReactionById(existingLikeId)
		if err != nil {
			return ferror.ReturnErr(err)
		}
	}
	return db.InsertDislikeToComment(user.Id, commentId)
}
