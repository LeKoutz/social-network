package db

import (
	"errors"
	"forum/src/utils"
)

type PostCategoryRow struct {
	PostId     int64
	CategoryId int64
}

func InsertPostCategory(pc PostCategoryRow) error {
	stmt, err := db.Prepare("INSERT INTO posts_categories (post_id, category_id) VALUES (?, ?)")
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	_, err = stmt.Exec(pc.PostId, pc.CategoryId)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}

func DeletePostCategoryByPostId(id int64) error {
	var err error
	_, err = db.Exec("DELETE FROM posts_categories WHERE post_id = ?", id)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}
