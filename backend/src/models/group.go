package models

import (
	"errors"
	"forum/src/db"
	"forum/src/ferror"
	"forum/src/utils"
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
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	if len(g.Title) >= 128 {
		err = ferror.ErrorGroupTitleTooLong
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}

// Adds a Group in the database. Returns its id or error
func (g *GroupType) Add() error {
	var err error
	err = g.ValidateGroup()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return g.InsertGroup()
}

func (g *GroupType) GetById() error {
	var err error
	err = g.SelectGroupById()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}

func (g *GroupType) IsMember(userId int64) (bool, error) {
	member, err := g.GroupRowType.IsMember(userId)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return false, err
	}
	return member, nil
}