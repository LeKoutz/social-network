package models

import (
	"errors"

	"forum/src/db"
	"forum/src/utils"
)

type CategoriesType []CategoryType

func (ct *CategoryType) FromCategoriesRowsType(crt db.CategoryRowType) {
	ct.CategoryRowType.Id = crt.Id
	ct.CategoryRowType.Description = crt.Description
	ct.CategoryRowType.Name = crt.Name
}

func (categories *CategoriesType) GetAll() error {
	var err error
	var categories_rows db.CategoriesRowsType
	categories_rows, err = db.GetAllCategories()
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	for _, item := range categories_rows {
		var category CategoryType
		category.FromCategoriesRowsType(item)
		*categories = append(*categories, category)
	}
	return nil
}

func (c *CategoriesType) IsEmpty() bool {
	if c == nil {
		return true
	}
	return len(*c) == 0
}

func (p *PostType) GetCategories() error {
	categories_rows, err := db.GetCategoriesByPostId(p.PostRowType.Id)
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return err
	}
	for _, item := range categories_rows {
		var category CategoryType
		category.FromCategoriesRowsType(item)
		p.Categories = append(p.Categories, category)
	}
	return nil
}
