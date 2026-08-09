package db

import (
	"database/sql"
	"errors"
	"forum/src/ferror"
	"forum/src/utils"
)

type PostRowType struct {
	Id              int64
	Title           string
	Body            string
	ImagePath       string
	UserId          int64
	TimestampString string
}

func (p *PostRowType) InsertPost() error {
	var query string = `
		INSERT INTO posts
		(title, body, image_path, user_id, timestamp)
		VALUES (?, ?, ?, ?, ?)
	`
	stmt, err := db.Prepare(query)
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	res, err := stmt.Exec(
		p.Title,
		p.Body,
		p.ImagePath,
		p.UserId,
		utils.GetCurrentTimestamp(),
	)
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	p.Id, err = res.LastInsertId()
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	return nil
}

func (p *PostRowType) UpdatePost() error {
	var err error
	_, err = db.Exec("UPDATE posts SET title = ?, body = ?, image_path = ? WHERE id = ?", p.Title, p.Body, p.ImagePath, p.Id)
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	return nil
}

func (p *PostRowType) SelectCommentsAndUsernameByPostId() (CommentRowsType, error) {
	rows, err := db.Query(`
	SELECT
	c.id,
	c.post_id,
	c.user_id,
	c.body,
	c.timestamp,
	u.username
	FROM comments c
	JOIN users u ON c.user_id = u.id
	WHERE c.post_id = ?
	ORDER BY c.timestamp ASC`, p.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return CommentRowsType{}, ferror.ErrorNoRows
		} else {
			return CommentRowsType{}, errors.Join(utils.GetFunctionName(), err)
		}
	}
	defer rows.Close()
	var comments CommentRowsType
	for rows.Next() {
		var comment CommentRowType
		err = rows.Scan(
			&comment.Id,
			&comment.PostId,
			&comment.UserId,
			&comment.Body,
			&comment.TimestampString,
			&comment.Username,
		)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return CommentRowsType{}, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (p *PostRowType) DeletePostById() error {
	tx, err := db.Begin()
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	_, err = tx.Exec("DELETE FROM reactions WHERE post_id = ?", p.Id)
	if err != nil {
		tx.Rollback()
		return errors.Join(utils.GetFunctionName(), err)
	}
	_, err = tx.Exec("DELETE FROM reactions WHERE comment_id IN (SELECT id FROM comments WHERE post_id = ?)", p.Id)
	if err != nil {
		tx.Rollback()
		return errors.Join(utils.GetFunctionName(), err)
	}
	_, err = tx.Exec("DELETE FROM comments WHERE post_id = ?", p.Id)
	if err != nil {
		tx.Rollback()
		return errors.Join(utils.GetFunctionName(), err)
	}
	_, err = tx.Exec("DELETE FROM posts_categories WHERE post_id = ?", p.Id)
	if err != nil {
		tx.Rollback()
		return errors.Join(utils.GetFunctionName(), err)
	}
	_, err = tx.Exec("DELETE FROM posts WHERE id = ?", p.Id)
	if err != nil {
		tx.Rollback()
		return errors.Join(utils.GetFunctionName(), err)
	}
	_, err = tx.Exec("DELETE FROM notifications WHERE post_id = ?", p.Id)
	if err != nil {
		tx.Rollback()
		return errors.Join(utils.GetFunctionName(), err)
	}
	return tx.Commit()
}

func (p *PostRowType) SelectPostById() error {
	var query string = `SELECT
		title, body, image_path, timestamp, user_id
		FROM posts
		WHERE id = ?`
	err := db.QueryRow(query, p.Id).Scan(&p.Title, &p.Body, &p.ImagePath, &p.TimestampString, &p.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Join(utils.GetFunctionName(), ferror.ErrorNoRows)
		} else {
			return errors.Join(utils.GetFunctionName(), err)
		}
	}
	// p.User, err = getUserById(p.UserId)
	// if err != nil {
	// 	err = errors.Join(utils.GetFunctionName(), err)
	// 	return errors.Join(utils.GetFunctionName(), err)
	// }
	return nil
}
