package forum

import (
	"slices"

	"golang.org/x/crypto/bcrypt"
)

func CompareRegistrationPasswords(pass1, pass2 string) bool {
	return pass1 == pass2
}

func IsUniqueUsername(username string) bool {
	usernames, err := getAllUsernames()
	if err != nil {
		var e Error
		e.Consume(err)
		e.LogError()
		return false
	}
	return !slices.Contains(usernames, username)
}

func IsUniqueEmail(email string) bool {
	emails, err := getAllUserEmails()
	if err != nil {
		var e Error
		e.Consume(err)
		e.LogError()
		return false
	}
	return !slices.Contains(emails, email)
}

func IsEmailRegistered(email string) bool {
	return !IsUniqueEmail(email)
}

func Auth(email, password string) error {
	var err error
	// I guess in order to authenticate against an email and a password, we will
	// need to first check if email is registered
	if !IsEmailRegistered(email) {
		// Return an error... this should be sent back to umh... places...?!
		return ErrorNotRegistered
	}
	var user User
	user, err = getUserByEmail(email)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

// Returns the hash (string) from password or error
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
