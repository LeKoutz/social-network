package models

import (
	"errors"
	"forum/src/db"
	"forum/src/utils"
)

type UsersType []UserType

func (u *UsersType) EditUsers() *UsersType {
	return u
}

func (u *UsersType) GetUsersForPanel(currentUserId int64) error {
	users, err := db.GetUsersForPanel(currentUserId)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	for _, user := range users {
		var x UserType
		x.UserRowType = user
		*u.EditUsers() = append(*u.EditUsers(), x)
	}
	return nil
}
