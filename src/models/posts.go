package models

import (
	"errors"
	"forum/src/db"
	"forum/src/utils"
)

type PostsType []PostType

func (p *PostsType) GetPostsByCategoryId(id int64) error {
	var err error
	var post_rows *db.PostRowsType = &db.PostRowsType{}
	err = post_rows.SelectPostsByCategoryId(id)
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return err
	}
	for _, i := range *post_rows {
		var post PostType
		post.FromPostRowType(&i)
		*p = append(*p, post)
	}
	return nil
}

func (p *PostsType) GetPosts() error {
	var err error
	var post_rows *db.PostRowsType = &db.PostRowsType{}
	err = post_rows.SelectAllPosts()
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return err
	}
	for _, i := range *post_rows {
		var post PostType
		post.FromPostRowType(&i)
		*p = append(*p, post)
	}
	return nil
}
