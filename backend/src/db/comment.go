package db

import (
	"forum/src/ferror"
	"forum/src/utils"
)

type CommentRowType struct {
	Id              int64
	PostId          int64
	UserId          int64
	Body            string
	ImagePath       string
	Timestamp       string
	Username        string
}

func (c *CommentRowType) InsertComment() error {
	res, err := db.Exec(
		"INSERT INTO comments (post_id, user_id, body, image_path, timestamp) VALUES (?, ?, ?, ?, ?)",
		c.PostId,
		c.UserId,
		c.Body,
		c.ImagePath,
		utils.GetCurrentTimestamp(),
	)
	c.Id, err = res.LastInsertId()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}

func (c *CommentRowType) SelectCommentById() error {
	err := db.QueryRow(
		`SELECT id, post_id, user_id, body, COALESCE(image_path, ''), timestamp
		FROM comments
		WHERE id = ?`, c.Id).Scan(&c.Id, &c.PostId, &c.UserId, &c.Body, &c.ImagePath, &c.Timestamp)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}

func (c *CommentRowType) DeleteCommentById() error {
	tx, err := db.Begin()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	_, err = tx.Exec("DELETE FROM reactions WHERE comment_id = ?", c.Id)
	if err != nil {
		tx.Rollback()
		return ferror.ReturnErr(err)
	}
	_, err = tx.Exec("DELETE FROM comments WHERE id = ?", c.Id)
	if err != nil {
		tx.Rollback()
		return ferror.ReturnErr(err)
	}
	_, err = tx.Exec("DELETE FROM notifications WHERE comment_id = ?", c.Id)
	if err != nil {
		tx.Rollback()
		return ferror.ReturnErr(err)
	}
	return tx.Commit()
}

func (c *CommentRowType) UpdateCommentById() error {
	_, err := db.Exec("UPDATE comments SET body = ?, image_path = ? WHERE id = ?", c.Body, c.ImagePath, c.Id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}
