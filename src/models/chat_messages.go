package models

import "forum/src/db"

type ChatMessagesType []ChatMessageType

func (m *ChatMessagesType) ConvertFromRow(t *db.ChatMessagesRowType) {
	var messages ChatMessagesType
	for _, i := range *t {
		var message ChatMessageType
		message.ConvertFromRow(&i)
		messages = append(messages, message)
	}
	m = &messages
}
