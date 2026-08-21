package db

import (
	"errors"
	"forum/src/utils"
)

type PostRowsType []PostRowType

func (posts *PostRowsType) SelectAllPosts() error {
	rows, err := db.Query(`SELECT id, title, body, timestamp, image_path FROM posts`)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var post PostRowType
		err = rows.Scan(&post.Id, &post.Title, &post.Body, &post.Timestamp, &post.ImagePath)
		if err != nil {
			if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return err
		}
		*posts = append(*posts, post)
	}
	return nil
}

func (posts *PostRowsType) SelectPostsByCategoryId(id int64) error {
	rows, err := db.Query(`
	SELECT posts.id, posts.title, posts.body, posts.timestamp, posts.image_path
	FROM posts
	JOIN posts_categories pc ON posts.id = pc.post_id
	JOIN categories ON pc.category_id = categories.id
	WHERE pc.category_id = ?`, id)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var post PostRowType
		err = rows.Scan(&post.Id, &post.Title, &post.Body, &post.Timestamp, &post.ImagePath)
		if err != nil {
			if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return err
		}
		*posts = append(*posts, post)
	}
	return nil
}
