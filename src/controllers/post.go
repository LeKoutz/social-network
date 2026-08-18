package controllers

import (
	"errors"
	"fmt"
	"forum/src/ferror"
	"forum/src/models"
	"forum/src/state"
	"forum/src/utils"
)

func GetPost(data state.StateController) error {
	var err error
	err = getPostDataById(data)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return data.EditUser().MarkAsReadPost(data.GetPost())
}

func CreatePost(data state.StateController) error {
	var err error
	err = data.EditPost().InsertPost()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	categories := data.EditPost().Categories
	for _, category := range categories {
		var post_cat models.PostCategory
		post_cat.PostCategoryRow.PostId = data.GetPost().Id
		post_cat.PostCategoryRow.CategoryId = category.Id
		err = post_cat.Add()
		if err != nil {
			if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return err
		}
	}
	data.SetRedirect(fmt.Sprintf("/post/view/%d", data.GetPost().Id))
	return nil
}

func getPostDataById(data state.StateController) error {
	var err error
	err = data.EditPost().SelectPostById()
	if err != nil {
		if err == ferror.ErrorNoRows {
			err = ferror.ErrorContentNotFound
			if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return err
		}
		return errors.Join(utils.GetFunctionName(), err)
	}
	data.EditPost().User.Id = data.GetPost().UserId
	err = data.EditPost().User.SelectUserById()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	err = data.EditPost().GetComments()
	if err != nil {
		if err == ferror.ErrorNoRows {
			err = ferror.ErrorContentNotFound
			if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return err
		}
		return errors.Join(utils.GetFunctionName(), err)
	}
	for i := range data.GetPost().Comments {
		err = data.EditPost().Comments[i].GetReactions()
		if err != nil {
			if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return err
		}
		err = data.EditPost().Comments[i].GetReactionsByUserId(data.GetUser().Id)
		if err != nil {
			if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return err
		}
	}
	err = data.EditPost().GetCategories()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	err = data.EditPost().GetReactions()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	err = data.EditPost().GetReactionsByUserId(data.GetUser().Id)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}

func ShowEditPost(data state.StateController) error {
	var err error
	err = data.EditPost().SelectPostById()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	if err = verifyUserPostAssociation(data); err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	err = data.EditCategories().GetAll()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	categories := data.GetCategories()
	err = data.EditPost().GetCategories()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	data.SetCategories(markSelectedCategories(categories, data.GetPost().Categories))
	data.SetEditPost(true)
	return nil
}

func verifyUserPostAssociation(data state.StateController) error {
	// Check your priviledge
	if data.GetPost().UserId != data.GetUser().Id {
		err := ferror.ErrorCommentPermissionDenied
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}

func UpdatePost(data state.StateController) error {
	var err error
	err = verifyUserPostAssociation(data)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	err = data.EditPost().Update()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	data.SetEditPost(false)
	return nil
}

func LikePost(data state.StateController) error {
	var err error
	err = data.EditUser().LikePost(data.GetPost().Id)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	utils.LogDebug(data.GetPost().Id)
	data.SetRedirect(fmt.Sprintf("/post/view/%d", data.GetPost().Id))
	return data.EditPost().CreateReactionNotification(data.GetUser().Id, "like")
}

func DislikePost(data state.StateController) error {
	var err error
	err = data.EditUser().DislikePost(data.GetPost().Id)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	data.SetRedirect(fmt.Sprintf("/post/view/%d", data.GetPost().Id))
	return data.EditPost().CreateReactionNotification(data.GetUser().Id, "dislike")
}

func PostReaction(data state.StateController) error {
	var err error
	err = data.EditPost().SelectPostById()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	switch data.GetRequest().FormValue("action") {
	case "like":
		return LikePost(data)
	case "dislike":
		return DislikePost(data)
	default:
		err = ferror.ErrorUnknownAction
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
}

func RemovePost(data state.StateController) error {
	var err error
	err = data.EditPost().SelectPostById()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	if data.GetPost().UserId != data.GetUser().Id {
		err = ferror.ErrorPostPermissionDenied
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	err = data.EditPost().Delete()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}
