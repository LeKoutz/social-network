package models

import (
	"forum/src/db"
	"forum/src/ferror"
)

type UsersType []UserType

func (u *UsersType) EditUsers() *UsersType {
	return u
}

func (u *UsersType) GetUsersForPanel(currentUserId int64) error {
	users, err := db.SelectUsersForPanel(currentUserId)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	for _, user := range users {
		var x UserType
		x.UserRowType = user
		*u.EditUsers() = append(*u.EditUsers(), x)
	}
	return nil
}
