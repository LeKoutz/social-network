package controllers

import (
	"forum/src/ferror"
	"forum/src/models"
	"forum/src/state"
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
	var activities models.ActivitiesType
	err = activities.GetActivityById(profile.UserId)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	profile.Activities = activities

	followers := &models.UserProfilesType{}
	tmp := &models.InvitationsType{}
	err = (*tmp).InvitationsRowType.SelectAcceptedInvitationsToUserId(profile.UserId)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	for _, req := range tmp.InvitationsRowType {
		var user models.UserProfileType
		user.UserId = req.FromUserId
		err = user.GetUserProfileIdentity()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		*followers = append(*followers, user)
	}
	profile.Followers = *followers

	following := &models.UserProfilesType{}
	tmp = &models.InvitationsType{}
	err = (*tmp).InvitationsRowType.SelectAcceptedInvitationsFromUserId(profile.UserId)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	for _, req := range tmp.InvitationsRowType {
		var user models.UserProfileType
		user.UserId = req.FromUserId
		err = user.GetUserProfileIdentity()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		*following = append(*following, user)
	}
	profile.Following = *following

	return nil
}
