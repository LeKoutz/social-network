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
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	err = data.EditPosts().GetPostsByCategoryId(data.GetCategory().Id)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	for i := range *data.EditPosts() {
		err = (*data.EditPosts())[i].GetReactions()
		if err != nil {
			if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return err
		}
		err = (*data.EditPosts())[i].GetReactionsByUserId(data.GetUser().Id)
		if err != nil {
			if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return err
		}
	}
	return nil
}
