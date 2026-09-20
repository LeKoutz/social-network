package controllers

import (
	"errors"
	"forum/src/state"
	"forum/src/utils"
)

func ShowGroups(data state.StateController) error {
	var err error
	err = data.EditGroups().GetGroups()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}
