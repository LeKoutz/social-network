package db

import (
	"errors"
	"forum/src/utils"
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
        if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
        return err
    }
    return nil
}
