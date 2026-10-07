package db

import (
	"forum/src/ferror"
)

type GroupMessagesRowType []GroupMessageRowType

func SelectGroupChatHistory(groupId, offset int64) (GroupMessagesRowType, []string, error) {
	rows, err := db.Query(`
      SELECT * FROM (
		SELECT gm.id, gm.group_id, gm.sender_id, gm.body, gm.timestamp, u.username
		FROM group_messages gm
		JOIN users u ON gm.sender_id = u.id
		WHERE gm.group_id = ?
		ORDER BY gm.id DESC
		LIMIT 10 OFFSET ?
      ) ORDER BY id ASC`, groupId, offset)
	if err != nil {
		return nil, nil, ferror.ReturnErr(err)
	}
	defer rows.Close()
	var messages GroupMessagesRowType
	var usernames []string
	for rows.Next() {
		var message GroupMessageRowType
		var username string
		err = rows.Scan(&message.Id, &message.GroupId, &message.SenderId, &message.Body, &message.Timestamp, &username)
		if err != nil {
			return nil, nil, ferror.ReturnErr(err)
		}
		messages = append(messages, message)
		usernames = append(usernames, username)
	}
	if err = rows.Err(); err != nil {
		return nil, nil, ferror.ReturnErr(err)
	}
	return messages, usernames, nil
}
