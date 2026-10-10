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
	groups := *data.EditGroups()
	for i := range groups {
		if (*data.EditGroups())[i].OwnerUserId == data.GetUser().Id {
			(*data.EditGroups())[i].Owner = true
		}
	}
	return nil
}
