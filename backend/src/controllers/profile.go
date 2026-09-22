package controllers

import (
	"forum/src/models"
	"forum/src/state"
	"forum/src/ferror"
)

func ShowProfile(data state.StateController) error {
	var err error
	profile := data.EditProfile()
	err = data.EditProfile().GetUserProfileIdentity()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	viewerId := data.GetUser().Id
	switch {
	case viewerId == profile.UserId:
		profile.CanView = true
	case !profile.PrivateProfile:
		profile.CanView = true
	default:
		// TODO: Check if users are following each other
		profile.CanView = false
	}
	if !profile.CanView {
		return nil
	}
	err = profile.GetUserProfileDetails()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	var postAuthor models.UserType
	postAuthor.Id = profile.UserId
	posts, err := postAuthor.GetPosts()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	data.SetPosts(posts)
	// TODO: Followers list for profile.UserId
	// TODO: Following list for profile.UserId

	return nil
}