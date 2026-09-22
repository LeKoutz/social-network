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
		return ferror.ReturnErr(ferror.ErrorPermissionDenied)
	}
	err = data.EditPosts().GetPostsByGroupId(data.GetGroup().Id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	posts := *data.EditPosts()
	for i := range posts {
		posts[i].User.Id = posts[i].UserId
		err = posts[i].User.GetById()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		err = posts[i].GetComments()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		for j := range posts[i].Comments {
			err = posts[i].Comments[j].GetReactions()
			if err != nil {
				return ferror.ReturnErr(err)
			}
			err = posts[i].Comments[j].GetReactionsByUserId(data.GetUser().Id)
			if err != nil {
				return ferror.ReturnErr(err)
			}
		}
		err = posts[i].GetCategories()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		err = posts[i].GetReactions()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		err = posts[i].GetReactionsByUserId(data.GetUser().Id)
		if err != nil {
			return ferror.ReturnErr(err)
		}
	}
	data.SetPosts(posts)
	return nil
}
