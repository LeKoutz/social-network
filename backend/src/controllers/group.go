package controllers

import (
	"errors"
	"forum/src/state"
	"forum/src/utils"
)

func CreateGroup(data state.StateController) error {
	var err error
	err = data.EditGroup().Add()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}