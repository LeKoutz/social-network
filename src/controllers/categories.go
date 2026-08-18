package controllers

import (
	"errors"
	"forum/src/state"
	"forum/src/utils"
)

func ShowCategories(data state.StateController) error {
	err := data.EditCategories().GetAll()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}
