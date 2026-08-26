package models

import "forum/src/db"

type ChatMessageType struct {
	db.ChatMessageRowType

	SenderUsername  string
}

func (m *ChatMessageType) Add() (int64, error) {
	return m.InsertMessage()
}

func (m *ChatMessageType) MarkAsRead() error {
	return m.UpdateMessageAsRead()
}
