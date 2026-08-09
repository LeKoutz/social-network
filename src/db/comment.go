package db

import (
	"errors"
	"forum/src/utils"
)

type CommentRowType struct {
	Id              int64
	PostId          int64
	UserId          int64
	Body            string
	TimestampString string
	Username        string
}

func (c *CommentRowType) InsertComment() error {
	res, err := db.Exec(
		"INSERT INTO comments (post_id, user_id, body, timestamp) VALUES (?, ?, ?, ?)",
		c.PostId,
		c.UserId,
		c.Body,
		utils.GetCurrentTimestamp(),
	)
	c.Id, err = res.LastInsertId()
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	return nil
}

func (c *CommentRowType) SelectCommentById() error {
	err := db.QueryRow(
		`SELECT id, post_id, user_id, body, timestamp
		FROM comments
		WHERE id = ?`, c.Id).Scan(&c.Id, &c.PostId, &c.UserId, &c.Body, &c.TimestampString)
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	return nil
}

func (c *CommentRowType) DeleteCommentById() error {
	tx, err := db.Begin()
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	_, err = tx.Exec("DELETE FROM reactions WHERE comment_id = ?", c.Id)
	if err != nil {
		tx.Rollback()
		return errors.Join(utils.GetFunctionName(), err)
	}
	_, err = tx.Exec("DELETE FROM comments WHERE id = ?", c.Id)
	if err != nil {
		tx.Rollback()
		return errors.Join(utils.GetFunctionName(), err)
	}
	_, err = tx.Exec("DELETE FROM notifications WHERE comment_id = ?", c.Id)
	if err != nil {
		tx.Rollback()
		return errors.Join(utils.GetFunctionName(), err)
	}
	return tx.Commit()
}

func (c *CommentRowType) UpdateCommentById() error {
	_, err := db.Exec("UPDATE comments SET body = ? WHERE id = ?", c.Body, c.Id)
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	return nil
}
