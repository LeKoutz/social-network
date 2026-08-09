package models

import "forum/src/db"

type ChatMessageType struct {
	Id              int64 `json:"Id"`
	SenderId        int64
	RecipientId     int64
	Body            string
	Timestamp       int64
	TimestampString string
	Read            bool
	SenderUsername  string
}

func (m *ChatMessageType) ConvertToRow() *db.ChatMessageRowType {
	t := &db.ChatMessageRowType{}
	t.Id = m.Id
	t.SenderId = m.SenderId
	t.RecipientId = m.RecipientId
	t.Body = m.Body
	t.Timestamp = m.Timestamp
	t.TimestampString = m.TimestampString
	t.Read = m.Read
	t.SenderUsername = m.SenderUsername
	return t
}

func (m *ChatMessageType) ConvertFromRow(t *db.ChatMessageRowType) {
	m.Id = t.Id
	m.SenderId = t.SenderId
	m.RecipientId = t.RecipientId
	m.Body = t.Body
	m.Timestamp = t.Timestamp
	m.TimestampString = t.TimestampString
	m.Read = t.Read
	m.SenderUsername = t.SenderUsername
}

func (m *ChatMessageType) Add() (int64, error) {
	return m.ConvertToRow().Add()
}

func (m *ChatMessageType) MarkAsRead() error {
	return m.ConvertToRow().MarkAsRead()
}
