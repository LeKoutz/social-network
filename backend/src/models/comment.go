package models

import (
	"forum/src/db"
	"forum/src/ferror"
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
		return ferror.ReturnErr(err)
	}
	if len(c.Body) > 1000 {
		err := ferror.ErrorCommentTooLong
		return ferror.ReturnErr(err)
	}
	return nil
}

func (c *CommentType) Add() error {
	var err error
	err = c.ValidateComment()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	err = c.InsertComment()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}

func (c *CommentType) GetReactions() error {
	var err error
	c.Likes, err = db.SelectLikesCountByCommentId(c.Id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	c.Dislikes, err = db.SelectDislikesCountByCommentId(c.Id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}

func (c *CommentType) GetReactionsByUserId(user_id int64) error {
	var err error
	c.Liked, err = HasUserLikedComment(user_id, c.Id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	c.Disliked, err = HasUserDislikedComment(user_id, c.Id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}

func (c *CommentType) GetById() error {
	var err error
	err = c.SelectCommentById()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}

func (c *CommentType) Update() error {
	var err error
	err = c.ValidateComment()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	err = c.UpdateCommentById()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}

func (c *CommentType) Delete() error {
	var err error
	err = c.DeleteCommentById()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}
