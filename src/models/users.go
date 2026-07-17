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