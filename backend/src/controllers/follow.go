package controllers

import "forum/src/state"

func CreateFollowInvitation(data state.StateController) error {
	return data.EditFollowRequest().AddFollowInvite()
}

func UnfollowUser(data state.StateController) error {
    return data.EditFollowRequest().Unfollow()
}