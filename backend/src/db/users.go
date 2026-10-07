package db

import "forum/src/ferror"

func SelectAllUsernames() ([]string, error) {
	rows, err := db.Query(`SELECT username FROM users`)
	if err != nil {
		return []string{}, ferror.ReturnErr(err)
	}
	defer rows.Close()
	var usernames []string
	for rows.Next() {
		var email string
		err = rows.Scan(&email)
		if err != nil {
			return []string{}, ferror.ReturnErr(err)
		}
		usernames = append(usernames, email)
	}
	return usernames, nil
}

func SelectAllUserEmails() ([]string, error) {
	rows, err := db.Query(`SELECT email FROM users`)
	if err != nil {
		return []string{}, ferror.ReturnErr(err)
	}
	defer rows.Close()
	var emails []string
	for rows.Next() {
		var email string
		err = rows.Scan(&email)
		if err != nil {
			return []string{}, ferror.ReturnErr(err)
		}
		emails = append(emails, email)
	}
	return emails, nil
}

func SelectAllUsers() ([]UserRowType, error) {
	rows, err := db.Query(`SELECT id, username FROM users`)
	if err != nil {
		return []UserRowType{}, ferror.ReturnErr(err)
	}
	defer rows.Close()
	var users []UserRowType
	for rows.Next() {
		var user UserRowType
		err = rows.Scan(&user.Id, &user.Username)
		if err != nil {
			return []UserRowType{}, ferror.ReturnErr(err)
		}
		users = append(users, user)
	}
	return users, nil
}

// Selects users that the currentUser can chat with (followers or following).
// It also retrieves the timestamp of the last message sent or received by each user
func SelectUsersForPanel(currentUserId int64) ([]UserRowType, error) {
	rows, err := db.Query(`
	SELECT
		u.id,
		u.username,
		COALESCE(MAX(CAST(m.timestamp AS INTEGER)), 0)
	FROM users u
	LEFT JOIN messages m ON
		(m.sender_id = ? AND m.recipient_id = u.id)
		OR
		(m.recipient_id = ? AND m.sender_id = u.id)
	WHERE u.id != ?
		AND EXISTS (
			SELECT 1 FROM invitations
			WHERE status = 'accepted'
			AND (
				(from_user_id = ? AND to_user_id = u.id)
				OR (from_user_id = u.id AND to_user_id = ?)
			)
		)
	GROUP BY u.id, u.username
	`, currentUserId, currentUserId, currentUserId, currentUserId, currentUserId)
	if err != nil {
		return []UserRowType{}, ferror.ReturnErr(err)
	}
	defer rows.Close()
	var users []UserRowType
	for rows.Next() {
		var user UserRowType
		err = rows.Scan(
			&user.Id,
			&user.Username,
			&user.LastMessageTimestamp, // TODO: Currently inconsistent. Should be in UserType
		)
		if err != nil {
			return []UserRowType{}, ferror.ReturnErr(err)
		}
		users = append(users, user)
	}
	return users, nil
}

func SelectUsersWithChats(currentUserId int64) ([]UserRowType, []ChatMessageRowType, error) {
	rows, err := db.Query(`
	SELECT
			u.id,
			u.username,
			lm.id,
			lm.sender_id,
			lm.recipient_id,
			lm.body,
			lm.timestamp,
			lm.read
	FROM users u
	JOIN messages lm ON lm.id = (
			SELECT m.id FROM messages m
			WHERE (m.sender_id = ? AND m.recipient_id = u.id)
					OR (m.recipient_id = ? AND m.sender_id = u.id)
			ORDER BY CAST(m.timestamp AS INTEGER) DESC, m.id DESC
			LIMIT 1
	)
	WHERE u.id != ?
	ORDER BY CAST(lm.timestamp AS INTEGER) DESC, lm.id DESC
	`, currentUserId, currentUserId, currentUserId)
	if err != nil {
		return []UserRowType{}, []ChatMessageRowType{}, ferror.ReturnErr(err)
	}
	defer rows.Close()
	var users []UserRowType
	var lastMessages []ChatMessageRowType
	for rows.Next() {
		var user UserRowType
		var message ChatMessageRowType
		err = rows.Scan(
			&user.Id,
			&user.Username,
			&message.Id,
			&message.SenderId,
			&message.RecipientId,
			&message.Body,
			&message.Timestamp,
			&message.Read,
		)
		if err != nil {
			return []UserRowType{}, []ChatMessageRowType{}, ferror.ReturnErr(err)
		}
		users = append(users, user)
		lastMessages = append(lastMessages, message)
	}
	if err = rows.Err(); err != nil {
		return []UserRowType{}, []ChatMessageRowType{}, ferror.ReturnErr(err)
	}
	return users, lastMessages, nil
}
