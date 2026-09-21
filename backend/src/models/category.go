package models

import (
	"forum/src/db"
	"forum/src/ferror"
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
		return ferror.ReturnErr(err)
	}
	if len(c.Name) >= 128 {
		err = ferror.ErrorCategoryNameTooLong
		return ferror.ReturnErr(err)
	}
	return nil
}

func (c *CategoryType) Add() error {
	var err error
	err = c.ValidateCategory()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return c.InsertCategory()
}

func (c *CategoryType) GetById() error {
	var err error
	err = c.SelectCategoryById()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}
