package db

import (
	"errors"
	"forum/src/utils"
)

type ChatMessageRowType struct {
	Id              int64 `json:"Id"`
	SenderId        int64
	RecipientId     int64
	Body            string
	Timestamp       int64
	TimestampString string
	Read            bool
	SenderUsername  string
}

func (msg *ChatMessageRowType) InsertMessage() (int64, error) {
	stmt, err := db.Prepare("INSERT INTO messages (sender_id, recipient_id, body, timestamp) VALUES (?, ?, ?, ?)")
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return 0, err
	}
	res, err := stmt.Exec(msg.SenderId, msg.RecipientId, msg.Body, msg.TimestampString)
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return 0, err
	}
	msgId, err := res.LastInsertId()
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return 0, err
	}
	return msgId, nil
}

func (msg *ChatMessageRowType) UpdateMessageAsRead() error {
	_, err := db.Exec(`UPDATE messages SET read = 1 WHERE id = ?`, msg.Id)
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return err
	}
	return nil
}
