package models

import (
	"errors"
	"forum/src/utils"
)

type Category struct {
	Id          int64
	Name        string
	Description string
	Selected    bool
}

func (c *Category) IsEmpty() bool {
	return c == nil || *c == Category{}
}

func GetCategoryById(id int64) (Category, error) {
	var category Category
	err := db.QueryRow(`SELECT name, description FROM categories WHERE id = ?`, id).Scan(&category.Name, &category.Description)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return Category{}, err
	}
	category.Id = id
	return category, nil
}

func (c *Category) ValidateCategory() error {
	if len((*c).Name) == 0 {
		return ErrorCategoryNameEmpty
	}
	if len((*c).Name) >= 128 {
		return ErrorCategoryNameTooLong
	}
	return nil
}

func (c *Category) DoesCategoryExist() bool {
	categories, err := GetAllCategories()
	if err != nil {
		(&Error{}).Consume(err).LogError()
		return false
	}
	for _, category := range categories {
		if category.Name == (*c).Name {
			return true
		}
	}
	return false
}
