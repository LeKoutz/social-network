package db

import (
	"database/sql"
	"errors"
	"forum/src/ferror"
)

type UserProfileRowType struct {
	UserId         int64
	Nickname       string
	FirstName      string
	LastName       string
	DateOfBirth    string
	About          string
	AvatarURL      string
	PrivateProfile bool
}

func (u *UserProfileRowType) InsertUserProfile() error {
    _, err := db.Exec(`
        INSERT INTO user_profiles (
            user_id, nickname, first_name, last_name, date_of_birth, about, avatar_url, private_profile
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `, u.UserId, u.Nickname, u.FirstName, u.LastName, u.DateOfBirth, u.About, u.AvatarURL, u.PrivateProfile)
    if err != nil {
        return ferror.ReturnErr(err)
    }
    return nil
}

func (u *UserProfileRowType) SelectUserProfileIdentity() (UserProfileRowType, string, error) {
	var profile UserProfileRowType
	var username string
	err := db.QueryRow(`
		SELECT up.user_id, u.username, COALESCE(up.avatar_url, ''), up.private_profile
		FROM user_profiles up
		JOIN users u ON u.id = up.user_id
		WHERE up.user_id = ?
	`, u.UserId).Scan(
		&profile.UserId, &username, &profile.AvatarURL, &profile.PrivateProfile,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ferror.ErrorNoRows
		}
		return profile, "", ferror.ReturnErr(err)
	}
	return profile, username, nil
}

func (u *UserProfileRowType) SelectUserProfileDetails() (UserProfileRowType, string, error) {
	var profile UserProfileRowType
	var email string
	err := db.QueryRow(`
		SELECT
		up.user_id, COALESCE(up.nickname, ''), up.first_name, up.last_name,
		COALESCE(up.date_of_birth, ''), COALESCE(up.about, ''), COALESCE(up.avatar_url, ''), up.private_profile,
		u.email
		FROM user_profiles up
		JOIN users u ON u.id = up.user_id
		WHERE up.user_id = ?
	`, u.UserId).Scan(
		&profile.UserId, &profile.Nickname, &profile.FirstName, &profile.LastName,
		&profile.DateOfBirth, &profile.About, &profile.AvatarURL, &profile.PrivateProfile,
		&email,
	)
	if err != nil {
		return profile, "", ferror.ReturnErr(err)
	}
	return profile, email, nil
}
