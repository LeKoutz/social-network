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

// GetUsersForPanel retrieves all users from the database, excluding the current user, and returns them as a slice of User structs.
// It also retrieves the timestamp of the last message sent or received by each user.
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
	GROUP BY u.id, u.username
	`, currentUserId, currentUserId, currentUserId)
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
