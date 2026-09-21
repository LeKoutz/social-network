package controllers

import (
	"forum/src/ferror"
	"forum/src/state"
)

func CreateGroup(data state.StateController) error {
	var err error
	err = data.EditGroup().Add()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}

func ShowGroup(data state.StateController) error {
	var err error
	err = data.EditGroup().GetById()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	member, err := data.EditGroup().IsMember(data.GetUser().Id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	data.EditGroup().Member = member
	if !member {
		return nil
	}
	err = data.EditPosts().GetPostsByGroupId(data.GetGroup().Id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	data.EditGroup().Posts = *data.EditPosts()
	for i := range data.EditGroup().Posts {
		data.EditGroup().Posts[i].User.Id = data.EditGroup().Posts[i].UserId
		err = data.EditGroup().Posts[i].User.GetById()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		err = data.EditGroup().Posts[i].GetComments()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		for j := range data.EditGroup().Posts[i].Comments {
			err = data.EditGroup().Posts[i].Comments[j].GetReactions()
			if err != nil {
				return ferror.ReturnErr(err)
			}
			err = data.EditGroup().Posts[i].Comments[j].GetReactionsByUserId(data.GetUser().Id)
			if err != nil {
				return ferror.ReturnErr(err)
			}
		}
		err = data.EditGroup().Posts[i].GetCategories()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		err = data.EditGroup().Posts[i].GetReactions()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		err = data.EditGroup().Posts[i].GetReactionsByUserId(data.GetUser().Id)
		if err != nil {
			return ferror.ReturnErr(err)
		}
	}
	return nil
}
