package db

import (
	"forum/src/ferror"
)

type GroupMessageRowType struct {
	Id        int64
	GroupId   int64
	SenderId  int64
	Body      string
	Timestamp string
}

func (msg *GroupMessageRowType) InsertGroupMessage() (int64, error) {
	res, err := db.Exec(
		`INSERT INTO group_messages (group_id, sender_id, body, timestamp) VALUES (?, ?, ?, ?)`,
		msg.GroupId, msg.SenderId, msg.Body, msg.Timestamp,
	)
	if err != nil {
		return 0, ferror.ReturnErr(err)
	}
	msgId, err := res.LastInsertId()
	if err != nil {
		return 0, ferror.ReturnErr(err)
	}
	return msgId, nil
}

func (msg *GroupMessageRowType) UpdateLastReadMessageId(userId int64) error {
	_, err := db.Exec(
		`UPDATE group_message_reads
		SET last_read_message_id = MAX(last_read_message_id, ?)
		WHERE group_id = ? AND user_id = ?
		AND EXISTS (SELECT 1 FROM group_messages WHERE id = ? AND group_id = ?)`,
		msg.Id, msg.GroupId, userId, msg.Id, msg.GroupId,
	)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}
