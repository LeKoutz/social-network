package db

import (
	"errors"
	"forum/src/utils"
)

type PostCategoryRow struct {
	PostId     int64
	CategoryId int64
}

func AddPostCategory(pc PostCategoryRow) error {
	stmt, err := db.Prepare("INSERT INTO posts_categories (post_id, category_id) VALUES (?, ?)")
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return errors.Join(utils.GetFunctionName(), err)
	}
	_, err = stmt.Exec(pc.PostId, pc.CategoryId)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return errors.Join(utils.GetFunctionName(), err)
	}
	return nil
}

func DeletePostCategoryByPostId(id int64) error {
	var err error
	_, err = db.Exec("DELETE FROM posts_categories WHERE post_id = ?", id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return errors.Join(utils.GetFunctionName(), err)
	}
	return nil
}
