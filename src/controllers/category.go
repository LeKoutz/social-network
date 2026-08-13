package controllers

import (
	"errors"
	"forum/src/state"
	"forum/src/utils"
)

func ShowCategory(data state.StateController) error {
	var err error
	err = data.EditCategory().SelectCategoryById()
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	err = data.EditPosts().GetPostsByCategoryId(data.GetCategory().Id)
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	for i := range *data.EditPosts() {
		err = (*data.EditPosts())[i].GetReactions()
		if err != nil {
			return errors.Join(utils.GetFunctionName(), err)
		}
		err = (*data.EditPosts())[i].GetReactionsByUserId(data.GetUser().Id)
		if err != nil {
			return errors.Join(utils.GetFunctionName(), err)
		}
	}
	return nil
}
