package models

import (
	"errors"
	"forum/src/db"
	"forum/src/ferror"
	"forum/src/utils"
)

type CategoryType struct {
	db.CategoryRowType
	Selected bool
}

func (c *CategoryType) IsEmpty() bool {
	return c == nil || *c == CategoryType{}
}

func (c *CategoryType) ValidateCategory() error {
	var err error
	if len(c.Name) == 0 {
		err = ferror.ErrorCategoryNameEmpty
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	if len(c.Name) >= 128 {
		err = ferror.ErrorCategoryNameTooLong
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}

func (c *CategoryType) Add() error {
	var err error
	err = c.ValidateCategory()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return c.InsertCategory()
}

func (c *CategoryType) GetById() error {
	var err error
	err = c.SelectCategoryById()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}
