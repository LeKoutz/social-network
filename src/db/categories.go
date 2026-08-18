package db

import (
	"errors"
	"forum/src/utils"
)

type CategoriesRowsType []CategoryRowType

func GetAllCategories() (CategoriesRowsType, error) {
	rows, err := db.Query(`SELECT id, name, description FROM categories`)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return nil, err
	}
	defer rows.Close()
	var categories CategoriesRowsType
	for rows.Next() {
		var category CategoryRowType
		err = rows.Scan(&category.Id, &category.Name, &category.Description)
		if err != nil {
			if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func GetCategoriesByPostId(post_id int64) (CategoriesRowsType, error) {
	var categories CategoriesRowsType
	var err error
	rows, err := db.Query(`
	SELECT c.id, c.name, c.description
	FROM categories c
	JOIN posts_categories pc ON c.id = pc.category_id
	WHERE pc.post_id = ?
	`, post_id)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return categories, err
	}
	defer rows.Close()
	for rows.Next() {
		var category CategoryRowType
		err = rows.Scan(&category.Id, &category.Name, &category.Description)
		if err != nil {
			if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return categories, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}
