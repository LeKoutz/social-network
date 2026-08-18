package models

import (
	"errors"
	"forum/src/utils"
	"forum/src/db"
)

type ChatMessagesType []ChatMessageType

func (m *ChatMessagesType) GetUnreadMessageIds(userId int64) error {
	var messages ChatMessagesType
	rows, err := db.SelectUnreadMessageIds(userId)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
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
	rows, err := db.SelectChatHistory(userId1, userId2, offset)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	for _, row := range rows {
		var message ChatMessageType
		message.ChatMessageRowType = row
		messages = append(messages, message)
	}
	*m = messages
	return nil
}
