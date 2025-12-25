package models

import (
	"errors"
	"forum/src/utils"
)

type Comment struct {
	Id              int64
	PostId          int64
	UserId          int64
	Body            string
	Timestamp       int64
	TimestampString string
	Likes           int64
	Liked           bool
	Dislikes        int64
	Disliked        bool
	Username        string
}

func (c *Comment) ValidateComment() error {
	if len(c.Body) == 0 {
		return ErrorCommentEmpty
	}
	if len(c.Body) > 1000 {
		return ErrorCommentTooLong
	}
	return nil
}

func AddComment(comment Comment) (int64, error) {
	if err := comment.ValidateComment(); err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return 0, err
	}
	res, err := DB.Exec(
		"INSERT INTO comments (post_id, user_id, body, timestamp) VALUES (?, ?, ?, ?)",
		comment.PostId,
		comment.UserId,
		comment.Body,
		utils.GetCurrentTimestamp(),
	)
	commentId, err := res.LastInsertId()
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return 0, err
	}
	return commentId, nil
}

func (c *Comment) GetReactions() error {
	var err error
	(*c).Likes, err = getLikesCountByCommentId((*c).Id)
	if err != nil {
		return err
	}
	(*c).Dislikes, err = getDislikesCountByCommentId((*c).Id)
	if err != nil {
		return err
	}
	return nil
}

func (c *Comment) GetReactionsByUserId(user_id int64) error {
	var err error
	(*c).Liked, err = HasUserLikedComment(user_id, (*c).Id)
	if err != nil {
		return err
	}
	(*c).Disliked, err = HasUserDislikedComment(user_id, (*c).Id)
	if err != nil {
		return err
	}
	return nil
}
