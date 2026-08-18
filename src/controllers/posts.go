package controllers

import (
	"errors"
	"forum/src/models"
	"forum/src/state"
	"forum/src/utils"
)

func markSelectedCategories(categories models.CategoriesType, selected models.CategoriesType) models.CategoriesType {
	selectedIDs := make(map[int64]bool, len(selected))
	for _, category := range selected {
		selectedIDs[category.Id] = true
	}
	for i := range categories {
		categories[i].Selected = selectedIDs[categories[i].Id]
	}
	return categories
}

func ShowPosts(data state.StateController) error {
	var err error
	err = data.EditPosts().GetPosts()
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
	// data.SetPosts(posts)
	return nil
}
