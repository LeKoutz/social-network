package controllers

import (
	"forum/src/ferror"
	"forum/src/state"
)

func ShowCategory(data state.StateController) error {
	var err error
	err = data.EditCategory().GetById()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	err = data.EditPosts().GetPostsByCategoryId(data.GetCategory().Id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	for i := range *data.EditPosts() {
		err = (*data.EditPosts())[i].GetReactions()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		err = (*data.EditPosts())[i].GetReactionsByUserId(data.GetUser().Id)
		if err != nil {
			return ferror.ReturnErr(err)
		}
	}
	return nil
}
