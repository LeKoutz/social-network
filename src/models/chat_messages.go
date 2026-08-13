package models

import (
	"errors"
	"forum/src/utils"
	"forum/src/db"
)

type ChatMessagesType []ChatMessageType

func (m *ChatMessagesType) convertFromRow(rows *db.ChatMessagesRowType) {
	var messages ChatMessagesType
	for _, i := range *rows {
		var message ChatMessageType
		message.ConvertFromRow(&i)
		messages = append(messages, message)
	}
	*m = messages
}

func (m *ChatMessagesType) GetUnreadMessageIds(userId int64) error {
	rows, err := db.SelectUnreadMessageIds(userId)
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return err
	}
	m.convertFromRow(&rows)
	return nil
}

func (m *ChatMessagesType) GetChatHistory(userId1, userId2, offset int64) error {
	rows, err := db.SelectChatHistory(userId1, userId2, offset)
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return err
	}
	m.convertFromRow(&rows)
	return nil
}
