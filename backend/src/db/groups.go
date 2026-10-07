package db

import "forum/src/ferror"

type GroupRowsType []GroupRowType

func (groups *GroupRowsType) SelectAllGroups() error {
	rows, err := db.Query(`
		SELECT groups.id, groups.timestamp, groups.title, groups.description, groups.owner_user_id, users.username
		FROM groups
		JOIN users ON groups.owner_user_id = users.id
	`)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	defer rows.Close()
	for rows.Next() {
		var group GroupRowType
		err = rows.Scan(
			&group.Id,
			&group.Timestamp,
			&group.Title,
			&group.Description,
			&group.OwnerUserId,
			&group.OwnerUsername,
		)
		if err != nil {
			return ferror.ReturnErr(err)
		}
		*groups = append(*groups, group)
	}
	return nil
}

func SelectGroupsWithUnreadMessages(userId int64) (GroupRowsType, GroupMessagesRowType, error) {
	rows, err := db.Query(`
		SELECT groups.id, groups.timestamp, groups.title, groups.description,
			groups.owner_user_id, users.username, COALESCE(gm.id, 0)
		FROM groups
		JOIN users ON groups.owner_user_id = users.id
		LEFT JOIN group_message_reads gmr
			ON gmr.group_id = groups.id AND gmr.user_id = ?
		LEFT JOIN group_messages gm
			ON gm.group_id = groups.id
			AND gm.id > gmr.last_read_message_id
			AND gm.sender_id != ?
		WHERE groups.owner_user_id = ?
		OR EXISTS (SELECT 1 FROM group_invitations gi
			WHERE gi.group_id = groups.id AND gi.to_user_id = ?
			AND gi.status = 'accepted')
		ORDER BY groups.id, gm.id
	`, userId, userId, userId, userId)
	if err != nil {
		return nil, nil, ferror.ReturnErr(err)
	}
	defer rows.Close()
	var groups GroupRowsType
	var messages GroupMessagesRowType
	for rows.Next() {
		var group GroupRowType
		var message GroupMessageRowType
		err = rows.Scan(
			&group.Id,
			&group.Timestamp,
			&group.Title,
			&group.Description,
			&group.OwnerUserId,
			&group.OwnerUsername,
			&message.Id,
		)
		if err != nil {
			return nil, nil, ferror.ReturnErr(err)
		}
		message.GroupId = group.Id
		groups = append(groups, group)
		messages = append(messages, message)
	}
	if err = rows.Err(); err != nil {
		return nil, nil, ferror.ReturnErr(err)
	}
	return groups, messages, nil
}
