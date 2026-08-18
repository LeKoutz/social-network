package models

import "forum/src/db"

type PostCategory struct {
	db.PostCategoryRow
}

func (pc *PostCategory) Add() error {
	return db.InsertPostCategory(pc.PostCategoryRow)
}
