package controllers

import (
	"errors"
	"forum/src/ferror"
	"forum/src/models"
	"forum/src/state"
	"forum/src/utils"
	"strings"
)

func ParsePostId(data state.StateController) (int64, error) {
	var postIdStr string
	var postId int64
	var ok bool
	var err error
	postIdStr = data.GetRequest().FormValue("post-id")
	if len(postIdStr) != 0 {
		utils.LogDebug(postIdStr)
		goto Convert
	}
	postIdStr, ok = strings.CutPrefix(data.GetRequest().RequestURI, "/api/post/view/")
	if ok && len(postIdStr) != 0 {
		utils.LogDebug(postIdStr)
		goto Convert
	}
	if len(postIdStr) == 0 {
		utils.LogDebug(postIdStr)
		return 0, ferror.ErrorPostEmptyId
	}
Convert:
	postId, err = utils.StringToInt64(postIdStr)
	if err != nil || postId == 0 {
		utils.LogDebug(postIdStr)
		return 0, ferror.ErrorInvalidPostId
	}
	return postId, nil
}

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
	data.EditPost().Id, err = ParsePostId(data)
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	err = data.EditPosts().GetPosts()
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
	// data.SetPosts(posts)
	return nil
}
