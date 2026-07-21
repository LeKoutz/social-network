package models

import (
	"errors"
	"forum/src/utils"
)

func GetAllUsernames() ([]string, error) {
	rows, err := db.Query(`SELECT username FROM users`)
	if err != nil {
		if config.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return []string{}, err
	}
	defer rows.Close()
	var usernames []string
	for rows.Next() {
		var email string
		err = rows.Scan(&email)
		if err != nil {
			if config.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return []string{}, err
		}
		usernames = append(usernames, email)
	}
	return usernames, nil
}

func GetAllUserEmails() ([]string, error) {
	rows, err := db.Query(`SELECT email FROM users`)
	if err != nil {
		if config.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return []string{}, err
	}
	defer rows.Close()
	var emails []string
	for rows.Next() {
		var email string
		err = rows.Scan(&email)
		if err != nil {
			if config.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return []string{}, err
		}
		emails = append(emails, email)
	}
	return emails, nil
}

func GetAllUsers() ([]User, error) {
	rows, err := db.Query(`SELECT id, username FROM users`)
	if err != nil {
		if config.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return []User{}, err
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.Id, &user.Username)
		if err != nil {
			if config.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return []User{}, err
		}
		users = append(users, user)
	}
	return users, nil
}

// GetUsersForPanel retrieves all users from the database, excluding the current user, and returns them as a slice of User structs.
// It also retrieves the timestamp of the last message sent or received by each user.
func GetUsersForPanel(currentUserId int64) ([]User, error) {
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
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return []User{}, err
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var user User
		err = rows.Scan(
			&user.Id,
			&user.Username,
			&user.LastMessageTimestamp,
		)
		if err != nil {
			if config.Debug {
				err = errors.Join(utils.GetFunctionName(), err)
			}
			return []User{}, err
		}
		users = append(users, user)
	}
	return users, nil
}
