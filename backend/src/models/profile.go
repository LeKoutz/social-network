package models

import (
	"forum/src/db"
	"forum/src/ferror"
	"mime/multipart"
	"time"
)

type UserProfileType struct {
	db.UserProfileRowType

	Username     string
	Email        string
	Activities   ActivitiesType
	CanView      bool
	Followers    UserProfilesType
	Following    UserProfilesType
	AvatarFile   multipart.File `json:"-"`
	RemoveAvatar bool           `json:"-"`
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

func (u *UserProfileType) ValidateProfile() error {
	if u.FirstName == "" || u.LastName == "" {
		return ferror.ReturnErr(ferror.ErrorBadRequest)
	}

	if u.DateOfBirth != "" {
		if _, err := time.Parse("2006-01-02", u.DateOfBirth); err != nil {
			return ferror.ReturnErr(ferror.ErrorBadRequest)
		}
		if u.DateOfBirth > time.Now().Format("2006-01-02") {
			return ferror.ReturnErr(ferror.ErrorBadRequest)
		}
	}

	return nil
}
