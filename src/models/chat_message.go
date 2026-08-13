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
	row := &db.ChatMessageRowType{}
	row.Id = m.Id
	row.SenderId = m.SenderId
	row.RecipientId = m.RecipientId
	row.Body = m.Body
	row.Timestamp = m.Timestamp
	row.TimestampString = m.TimestampString
	row.Read = m.Read
	row.SenderUsername = m.SenderUsername
	return row
}

func (m *ChatMessageType) ConvertFromRow(row *db.ChatMessageRowType) {
	m.Id = row.Id
	m.SenderId = row.SenderId
	m.RecipientId = row.RecipientId
	m.Body = row.Body
	m.Timestamp = row.Timestamp
	m.TimestampString = row.TimestampString
	m.Read = row.Read
	m.SenderUsername = row.SenderUsername
}

func (m *ChatMessageType) Add() (int64, error) {
	return m.ConvertToRow().Add()
}

func (m *ChatMessageType) MarkAsRead() error {
	return m.ConvertToRow().MarkAsRead()
}
