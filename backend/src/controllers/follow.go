package controllers

import (
	"forum/src/ferror"
	"forum/src/state"
)

func CreateFollowInvitation(data state.StateController) error {
	return data.EditFollowRequest().AddFollowInvite()
}

func UnfollowUser(data state.StateController) error {
	return data.EditFollowRequest().Unfollow()
}

func GetFollowsForUser(data state.StateController) error {
	err := data.EditFollows().GetForUser(data.GetUser().Id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}

func AcceptFollowRequest(data state.StateController) error {
	err := data.EditFollowRequest().Accept(data.GetUser().Id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}

func DenyFollowRequest(data state.StateController) error {
	err := data.EditFollowRequest().Deny(data.GetUser().Id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}
