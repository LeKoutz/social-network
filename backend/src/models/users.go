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

func (u *UsersType) GetAllChats(currentUserId int64) error {
	users, lastMessages, err := db.SelectUsersWithChats(currentUserId)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	for i, user := range users {
		var x UserType
		x.UserRowType = user

		var msg ChatMessageType
		msg.ChatMessageRowType = lastMessages[i]
		x.ChatMessages = ChatMessagesType{msg}

		*u = append(*u, x)
	}
	return nil
}
