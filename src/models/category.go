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
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), ferror.ErrorCategoryNameEmpty)
		}
		return err
	}
	if len(c.Name) >= 128 {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), ferror.ErrorCategoryNameTooLong)
		}
		return err
	}
	return nil
}

func (c *CategoryType) Add() error {
	var err error
	err = c.ValidateCategory()
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	return c.InsertCategory()
}
