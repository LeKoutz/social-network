package controllers

import (
	"forum/src/ferror"
	"forum/src/models"
	"forum/src/state"
)

func GetPost(data state.StateController) error {
	var err error
	err = getPostDataById(data)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return data.EditUser().MarkAsReadPost(data.GetPost())
}

func CreatePost(data state.StateController) error {
	var err error
	if data.GetPost().GroupId != 0 {
		err = verifyPostGroupAccess(data)
		if err != nil {
			return ferror.ReturnErr(err)
		}
	}
	err = data.EditPost().InsertPost()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	categories := data.EditPost().Categories
	for _, category := range categories {
		var post_cat models.PostCategory
		post_cat.PostCategoryRow.PostId = data.GetPost().Id
		post_cat.PostCategoryRow.CategoryId = category.Id
		err = post_cat.Add()
		if err != nil {
			return ferror.ReturnErr(err)
		}
	}
	return nil
}

func getPostDataById(data state.StateController) error {
	var err error
	err = data.EditPost().GetById()
	if err != nil {
		if err == ferror.ErrorNoRows {
			err = ferror.ErrorContentNotFound
			return ferror.ReturnErr(err)
		}
		return ferror.ReturnErr(err)
	}
	err = verifyPostGroupAccess(data)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	data.EditPost().User.Id = data.GetPost().UserId
	err = data.EditPost().User.GetById()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	err = data.EditPost().GetComments()
	if err != nil {
		if err == ferror.ErrorNoRows {
			err = ferror.ErrorContentNotFound
			return ferror.ReturnErr(err)
		}
		return ferror.ReturnErr(err)
	}
	for i := range data.GetPost().Comments {
		err = data.EditPost().Comments[i].GetReactions()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		err = data.EditPost().Comments[i].GetReactionsByUserId(data.GetUser().Id)
		if err != nil {
			return ferror.ReturnErr(err)
		}
	}
	err = data.EditPost().GetCategories()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	err = data.EditPost().GetReactions()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	err = data.EditPost().GetReactionsByUserId(data.GetUser().Id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}

func ShowEditPost(data state.StateController) error {
	var err error
	err = data.EditPost().GetById()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	err = verifyPostGroupAccess(data)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	if err = verifyUserPostAssociation(data); err != nil {
		return ferror.ReturnErr(err)
	}
	err = data.EditCategories().GetAll()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	categories := data.GetCategories()
	err = data.EditPost().GetCategories()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	data.SetCategories(markSelectedCategories(categories, data.GetPost().Categories))
	data.SetEditPost(true)
	return nil
}

func verifyUserPostAssociation(data state.StateController) error {
	// Check your priviledge
	if data.GetPost().UserId != data.GetUser().Id {
		err := ferror.ErrorCommentPermissionDenied
		return ferror.ReturnErr(err)
	}
	return nil
}

// verifyPostGroupAccess restricts group posts to the group's members.
func verifyPostGroupAccess(data state.StateController) error {
	var err error
	if data.GetPost().GroupId == 0 {
		return nil
	}
	data.EditGroup().Id = data.GetPost().GroupId
	member, err := data.EditGroup().IsMember(data.GetUser().Id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	if !member {
		err = ferror.ErrorGroupMembershipRequired
		return ferror.ReturnErr(err)
	}
	return nil
}

func UpdatePost(data state.StateController) error {
	var err error
	err = verifyUserPostAssociation(data)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	err = data.EditPost().Update()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	data.SetEditPost(false)
	return nil
}

func LikePost(data state.StateController) error {
	var err error
	err = data.EditUser().LikePost(data.GetPost().Id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return data.EditPost().CreateReactionNotification(data.GetUser().Id, "like")
}

func DislikePost(data state.StateController) error {
	var err error
	err = data.EditUser().DislikePost(data.GetPost().Id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return data.EditPost().CreateReactionNotification(data.GetUser().Id, "dislike")
}

func PostReaction(data state.StateController) error {
	var err error
	err = data.EditPost().GetById()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	err = verifyPostGroupAccess(data)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	switch data.GetRequest().FormValue("action") {
	case "like":
		return LikePost(data)
	case "dislike":
		return DislikePost(data)
	default:
		err = ferror.ErrorUnknownAction
		return ferror.ReturnErr(err)
	}
}

func RemovePost(data state.StateController) error {
	var err error
	err = data.EditPost().GetById()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	err = verifyPostGroupAccess(data)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	if data.GetPost().UserId != data.GetUser().Id {
		err = ferror.ErrorPostPermissionDenied
		return ferror.ReturnErr(err)
	}
	err = data.EditPost().Delete()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}
