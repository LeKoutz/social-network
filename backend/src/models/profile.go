package models

import (
	"forum/src/db"
	"forum/src/ferror"
)

type UserProfileType struct {
	db.UserProfileRowType

	Username             string
	Email                string
	Activities			 ActivitiesType
	CanView				 bool
}

func (u *UserProfileType) GetUserProfileIdentity() error {
	var err error
	profile, username, err := u.SelectUserProfileIdentity()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	u.UserProfileRowType = profile
	u.Username = username
	return nil
}

func (u *UserProfileType) GetUserProfileDetails() error {
	var err error
	profile, email, err := u.SelectUserProfileDetails()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	u.UserProfileRowType = profile
	u.Email = email
	return nil
}