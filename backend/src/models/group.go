package models

import (
	"forum/src/db"
	"forum/src/ferror"
)

type GroupType struct {
	db.GroupRowType

	Member bool
	Posts  PostsType
}

func (g *GroupType) ValidateGroup() error {
	var err error
	if len(g.Title) == 0 {
		err = ferror.ErrorGroupTitleEmpty
		return ferror.ReturnErr(err)
	}
	if len(g.Title) >= 128 {
		err = ferror.ErrorGroupTitleTooLong
		return ferror.ReturnErr(err)
	}
	return nil
}

// Adds a Group in the database. Returns its id or error
func (g *GroupType) Add() error {
	var err error
	err = g.ValidateGroup()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return g.InsertGroup()
}

func (g *GroupType) GetById() error {
	var err error
	err = g.SelectGroupById()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}

func (g *GroupType) IsMember(userId int64) (bool, error) {
	member, err := g.GroupRowType.IsMember(userId)
	if err != nil {
		return false, ferror.ReturnErr(err)
	}
	return member, nil
}
