package models

import (
	"errors"
	"forum/src/db"
	"forum/src/ferror"
	"forum/src/utils"
)

type CommentType struct {
	db.CommentRowType
	Likes     int64
	Liked     bool
	Dislikes  int64
	Disliked  bool
}

func (c *CommentType) ValidateComment() error {
	if len(c.Body) == 0 {
		err := ferror.ErrorCommentEmpty
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	if len(c.Body) > 1000 {
		err := ferror.ErrorCommentTooLong
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}

func (c *CommentType) Add() error {
	var err error
	err = c.ValidateComment()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	err = c.InsertComment()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}

func (c *CommentType) GetReactions() error {
	var err error
	c.Likes, err = db.SelectLikesCountByCommentId(c.Id)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	c.Dislikes, err = db.SelectDislikesCountByCommentId(c.Id)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}

func (c *CommentType) GetReactionsByUserId(user_id int64) error {
	var err error
	c.Liked, err = HasUserLikedComment(user_id, c.Id)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	c.Disliked, err = HasUserDislikedComment(user_id, c.Id)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}

func (c *CommentType) GetById() error {
	var err error
	err = c.SelectCommentById()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}

func (c *CommentType) Update() error {
	var err error
	err = c.ValidateComment()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	err = c.UpdateCommentById()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}

func (c *CommentType) Delete() error {
	var err error
	err = c.DeleteCommentById()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}
