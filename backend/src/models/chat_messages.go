package models

import (
	"forum/src/db"
	"forum/src/ferror"
)

type ChatMessagesType []ChatMessageType

func (m *ChatMessagesType) GetUnreadMessageIds(userId int64) error {
	var messages ChatMessagesType
	rows, err := db.SelectUnreadMessageIds(userId)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	for _, row := range rows {
		var message ChatMessageType
		message.ChatMessageRowType = row
		messages = append(messages, message)
	}
	*m = messages
	return nil
}

func (m *ChatMessagesType) GetChatHistory(userId1, userId2, offset int64) error {
	var messages ChatMessagesType
	rows, usernames, err := db.SelectChatHistory(userId1, userId2, offset)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	for i, row := range rows {
		var message ChatMessageType
		message.ChatMessageRowType = row
		message.SenderUsername = usernames[i]
		messages = append(messages, message)
	}
	*m = messages
	return nil
}
