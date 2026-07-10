package models

import (
	"errors"
	"forum/src/utils"
)

type ChatMessage struct {
	Id				int64
	SenderId		int64
	RecipientId		int64
	Body			string
	Timestamp		int64
	TimestampString	string
	Read			bool
	SenderUsername	string
}

func (msg *ChatMessage) Add() (error) {
	stmt, err := db.Prepare("INSERT INTO messages (sender_id, recipient_id, body, timestamp) VALUES (?, ?, ?, ?)")
	if err != nil {
		if config.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	_, err = stmt.Exec(msg.SenderId, msg.RecipientId, msg.Body, utils.GetCurrentTimestamp())
	if err != nil {
		if config.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}