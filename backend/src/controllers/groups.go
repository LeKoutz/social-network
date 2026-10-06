package controllers

import (
	"forum/src/ferror"
	"forum/src/state"
)

func ShowGroups(data state.StateController) error {
	var err error
	err = data.EditGroups().GetGroups()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}
