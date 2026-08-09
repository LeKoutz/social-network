package controllers

import (
	"errors"
	"forum/src/ferror"
	"forum/src/state"
	"forum/src/utils"
	"strings"
)

func ParseCategoryId(data state.StateController) (int64, error) {
	id, ok := strings.CutPrefix(data.GetRequest().RequestURI, "/api/category/view/")
	if !ok || len(id) == 0 {
		return 0, ferror.ErrorCategoryEmptyId
	}
	categoryId, err := utils.StringToInt64(id)
	if err != nil {
		return 0, ferror.ErrorInvalidCategoryId
	}
	return categoryId, nil
}

func ShowCategories(data state.StateController) error {
	err := data.EditCategories().GetAll()
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	return nil
}
