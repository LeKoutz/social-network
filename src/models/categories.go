package models

import (
	"errors"

	"forum/src/db"
	"forum/src/utils"
)

type CategoriesType []CategoryType

func (c *CategoriesType) GetAll() error {
	var categories CategoriesType
	var err error
	rows, err := db.SelectAllCategories()
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	for _, row := range rows {
		var category CategoryType
		category.CategoryRowType = row
		categories = append(categories, category)
	}
	*c = categories
	return nil
}

func (c *CategoriesType) IsEmpty() bool {
	if c == nil {
		return true
	}
	return len(*c) == 0
}

func (p *PostType) GetCategories() error {
	var categories CategoriesType
	rows, err := db.SelectCategoriesByPostId(p.PostRowType.Id)
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return err
	}
	for _, row := range rows {
		var category CategoryType
		category.CategoryRowType = row
		categories = append(categories, category)
	}
	p.Categories = categories
	return nil
}
