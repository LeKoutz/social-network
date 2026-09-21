package controllers

import (
	"forum/src/ferror"
	"forum/src/state"
)

func ShowCategories(data state.StateController) error {
	err := data.EditCategories().GetAll()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}
