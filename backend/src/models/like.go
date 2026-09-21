package models

import (
	"forum/src/db"
	"forum/src/ferror"
)

func (user *UserType) LikeComment(commentId int64) error {
	likeId , err := db.SelectUserLikeFromComment(user.Id, commentId)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	if likeId != 0 {
		return db.DeleteReactionById(likeId)
	}
	existingDislikeId, err := db.SelectUserDislikeFromComment(user.Id, commentId)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	if existingDislikeId != 0 {
		err = db.DeleteReactionById(existingDislikeId)
		if err != nil {
			return ferror.ReturnErr(err)
		}
	}
	return db.InsertLikeToComment(user.Id, commentId)
}

func (user *UserType) LikePost(postId int64) error {
	likeId, err := db.SelectUserLikeFromPost(user.Id, postId)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	if likeId != 0 {
		return db.DeleteReactionById(likeId)
	}
	existingDislikeId, err := db.SelectUserDislikeFromPost(user.Id, postId)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	if existingDislikeId != 0 {
		err = db.DeleteReactionById(existingDislikeId)
		if err != nil {
			return ferror.ReturnErr(err)
		}
	}
	return db.InsertLikeToPost(user.Id, postId)
}
