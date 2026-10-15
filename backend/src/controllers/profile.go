package controllers

import (
	"forum/src/ferror"
	"forum/src/models"
	"forum/src/state"
	"os"
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
		user.UserId = req.ToUserId
		err = user.GetUserProfileIdentity()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		*following = append(*following, user)
	}
	profile.Following = *following

	return nil
}

func ShowEditProfile(data state.StateController) error {
	if !data.GetUser().LoggedIn {
		return ferror.ReturnErr(ferror.ErrorUnauthorizedAction)
	}

	profile := data.EditProfile()
	profile.UserId = data.GetUser().Id
	profile.Username = data.GetUser().Username

	return profile.GetUserProfileDetails()
}

func UpdateProfile(data state.StateController) error {
	profile := data.EditProfile()

	if !data.GetUser().LoggedIn || profile.UserId != data.GetUser().Id {
		return ferror.ReturnErr(ferror.ErrorUnauthorizedAction)
	}

	if err := profile.ValidateProfile(); err != nil {
		return ferror.ReturnErr(err)
	}

	if profile.RemoveAvatar && profile.AvatarFile != nil {
		return ferror.ReturnErr(ferror.ErrorBadRequest)
	}

	previousAvatar := profile.AvatarURL
	var uploadedAvatar string

	if profile.AvatarFile != nil {
		var err error
		uploadedAvatar, err = models.SaveImage(profile.AvatarFile)
		if err != nil {
			return ferror.ReturnErr(err)
		}
		profile.AvatarURL = uploadedAvatar
	} else if profile.RemoveAvatar {
		profile.AvatarURL = ""
	}

	if err := profile.UpdateUserProfile(); err != nil {
		profile.AvatarURL = previousAvatar
		if uploadedAvatar != "" {
			_ = os.Remove(uploadedAvatar)
		}
		return ferror.ReturnErr(err)
	}

	data.SetMessage(models.Message{
		Has:     true,
		Type:    "Success",
		Content: "Profile updated successfully",
	})

	return nil
}
