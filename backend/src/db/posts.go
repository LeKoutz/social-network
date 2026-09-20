package db

import (
	"errors"
	"forum/src/utils"
)

type PostRowsType []PostRowType

func (posts *PostRowsType) SelectAllPosts() error {
	rows, err := db.Query(`SELECT id, title, body, timestamp, image_path FROM posts WHERE group_id IS NULL`)
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
	WHERE pc.category_id = ? AND posts.group_id IS NULL`, id)
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

func (posts *PostRowsType) SelectGroupPostsByGroupId(id int64) error {
	rows, err := db.Query(`
	SELECT posts.id, posts.title, posts.body, posts.timestamp, posts.image_path, posts.user_id,
	posts.group_id, groups.title
	FROM posts
	JOIN groups ON posts.group_id = groups.id
	WHERE posts.group_id = ?
	ORDER BY posts.timestamp DESC`, id)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var post PostRowType
		err = rows.Scan(
			&post.Id,
			&post.Title,
			&post.Body,
			&post.Timestamp,
			&post.ImagePath,
			&post.UserId,
			&post.GroupId,
			&post.GroupTitle,
		)
		if err != nil {
			if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return err
		}
		*posts = append(*posts, post)
	}
	return nil
}
