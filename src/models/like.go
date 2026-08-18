package models

import (
	"errors"
	"forum/src/db"
	"forum/src/utils"
)

func (user *UserType) LikeComment(commentId int64) error {
	likeId , err := db.SelectUserLikeFromComment(user.Id, commentId)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	if likeId != 0 {
		return db.DeleteReactionById(likeId)
	}
	existingDislikeId, err := db.SelectUserDislikeFromComment(user.Id, commentId)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	if existingDislikeId != 0 {
		err = db.DeleteReactionById(existingDislikeId)
		if err != nil {
			if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return err
		}
	}
	return db.InsertLikeToComment(user.Id, commentId)
}

func (user *UserType) LikePost(postId int64) error {
	likeId, err := db.SelectUserLikeFromPost(user.Id, postId)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	if likeId != 0 {
		return db.DeleteReactionById(likeId)
	}
	existingDislikeId, err := db.SelectUserDislikeFromPost(user.Id, postId)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	if existingDislikeId != 0 {
		err = db.DeleteReactionById(existingDislikeId)
		if err != nil {
			if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return err
		}
	}
	return db.InsertLikeToPost(user.Id, postId)
}
