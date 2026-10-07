package models

import (
	"forum/src/db"
	"forum/src/ferror"
)

type GroupMessagesType []GroupMessageType

func (m *GroupMessagesType) GetChatHistory(groupId, offset int64) error {
	var messages GroupMessagesType
	rows, usernames, err := db.SelectGroupChatHistory(groupId, offset)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	for i, row := range rows {
		var message GroupMessageType
		message.GroupMessageRowType = row
		message.SenderUsername = usernames[i]
		messages = append(messages, message)
	}
	*m = messages
	return nil
}
