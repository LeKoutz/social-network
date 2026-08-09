package controllers

import "forum/src/state"

func Index(data state.StateController) error {
	if data.GetUser().LoggedIn {
		return data.EditCategories().GetAll()
	}
	data.WriteResponse()
	return nil
}
