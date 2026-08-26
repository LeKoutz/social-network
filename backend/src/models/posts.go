package models

import (
	"errors"
	"forum/src/db"
	"forum/src/utils"
)

type PostsType []PostType

func (p *PostsType) GetPostsByCategoryId(id int64) error {
	var err error
	var posts PostsType
	var rows db.PostRowsType
	err = rows.SelectPostsByCategoryId(id)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	for _, row := range rows {
		var post PostType
		post.PostRowType = row
		posts = append(posts, post)
	}
	*p = posts
	return nil
}

func (p *PostsType) GetPosts() error {
	var err error
	var posts PostsType
	var rows db.PostRowsType
	err = rows.SelectAllPosts()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	for _, row := range rows {
		var post PostType
		post.PostRowType = row
		posts = append(posts, post)
	}
	*p = posts
	return nil
}
