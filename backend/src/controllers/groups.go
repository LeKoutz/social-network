package controllers

import (
	"forum/src/state"
	"forum/src/ferror"
)

func ShowGroups(data state.StateController) error {
	var err error
	err = data.EditGroups().GetGroups()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}
