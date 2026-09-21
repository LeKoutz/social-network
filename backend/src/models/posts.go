package models

import (
	"forum/src/db"
	"forum/src/ferror"
)

type PostsType []PostType

func (p *PostsType) GetPostsByCategoryId(id int64) error {
	var err error
	var posts PostsType
	var rows db.PostRowsType
	err = rows.SelectPostsByCategoryId(id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	for _, row := range rows {
		var post PostType
		post.PostRowType = row
		posts = append(posts, post)
	}
	*p = posts
	return nil
}

func (p *PostsType) GetPostsByGroupId(id int64) error {
	var err error
	var posts PostsType
	var rows db.PostRowsType
	err = rows.SelectGroupPostsByGroupId(id)
	if err != nil {
		return ferror.ReturnErr(err)
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
		return ferror.ReturnErr(err)
	}
	for _, row := range rows {
		var post PostType
		post.PostRowType = row
		posts = append(posts, post)
	}
	*p = posts
	return nil
}
