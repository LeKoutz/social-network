package controllers

import (
	"errors"
	"forum/src/ferror"
	"forum/src/state"
	"forum/src/utils"
	"strings"
)

func ShowCategories(data state.StateController) error {
	err := data.EditCategories().GetAll()
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	return nil
}
