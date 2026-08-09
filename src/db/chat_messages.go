package db

import (
	"errors"
	"forum/src/utils"
)

type ChatMessagesRowType []ChatMessageRowType

func GetChatHistory(userId1, userId2, offset int64) (ChatMessagesRowType, error) {
	rows, err := db.Query(`
	SELECT * FROM (
    SELECT m.id, m.sender_id, m.recipient_id, m.body, m.timestamp, u.username
    FROM messages m
    JOIN users u ON m.sender_id = u.id
    WHERE (m.sender_id = ? AND m.recipient_id = ?)
       OR (m.sender_id = ? AND m.recipient_id = ?)
    ORDER BY m.timestamp DESC
    LIMIT 10 OFFSET ?
) ORDER BY timestamp ASC`, userId1, userId2, userId2, userId1, offset)
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return ChatMessagesRowType{}, err
	}
	defer rows.Close()
	var messages ChatMessagesRowType
	for rows.Next() {
		var message ChatMessageRowType
		var ts string
		err = rows.Scan(&message.Id, &message.SenderId, &message.RecipientId, &message.Body, &ts, &message.SenderUsername)
		if err != nil {
			if config.Debug {
				err = errors.Join(utils.GetFunctionName(), err)
			}
			return ChatMessagesRowType{}, err
		}
		t, err := utils.ConvertStringToTime(ts)
		if err != nil {
			if config.Debug {
				err = errors.Join(utils.GetFunctionName(), err)
			}
			return ChatMessagesRowType{}, err
		}
		message.TimestampString = utils.ConvertTimeToString(t)
		messages = append(messages, message)
	}
	return messages, nil
}

func GetUnreadMessageIds(userId int64) (ChatMessagesRowType, error) {
	rows, err := db.Query(`
	SELECT id, sender_id
	FROM messages
	WHERE recipient_id = ?
	AND read = 0`, userId)
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return ChatMessagesRowType{}, err
	}
	defer rows.Close()
	var messages ChatMessagesRowType
	for rows.Next() {
		var message ChatMessageRowType
		err = rows.Scan(&message.Id, &message.SenderId)
		if err != nil {
			if config.Debug {
				err = errors.Join(utils.GetFunctionName(), err)
			}
			return ChatMessagesRowType{}, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}
