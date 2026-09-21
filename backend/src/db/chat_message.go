package db

import (
	"forum/src/ferror"
)

type ChatMessageRowType struct {
	Id              int64 `json:"Id"`
	SenderId        int64
	RecipientId     int64
	Body            string
	Timestamp       string
	Read            bool
}

func (msg *ChatMessageRowType) InsertMessage() (int64, error) {
	stmt, err := db.Prepare("INSERT INTO messages (sender_id, recipient_id, body, timestamp) VALUES (?, ?, ?, ?)")
	if err != nil {
		return 0, ferror.ReturnErr(err)
	}
	res, err := stmt.Exec(msg.SenderId, msg.RecipientId, msg.Body, msg.Timestamp)
	if err != nil {
		return 0, ferror.ReturnErr(err)
	}
	msgId, err := res.LastInsertId()
	if err != nil {
		return 0, ferror.ReturnErr(err)
	}
	return msgId, nil
}

func (msg *ChatMessageRowType) UpdateMessageAsRead() error {
	_, err := db.Exec(`UPDATE messages SET read = 1 WHERE id = ?`, msg.Id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}
