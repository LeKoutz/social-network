package forum

import (
	"crypto/rand"
	"slices"

	"golang.org/x/crypto/bcrypt"
)

func SaltGenerator(length int) string {
	return rand.Text()[:length]
}

func CompareRegistrationPasswords(pass1, pass2 string) bool {
	return pass1 == pass2
}

func IsUniqueUsername(username string) bool {
	// Check against all usernames?!
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
	// I guess in order to authenticate against an email and a password, we will
	// need to first check if email is registered
	if !IsEmailRegistered(email) {
		// Return an error... this should be sent back to umh... places...?!
		return ErrorNotRegistered
	}
	// Since we passed the initial test, we can fetch the user's salt
	salt := GetUserSalt(email)
	hashStored := GetUserHash(email)
	// Salt password
	saltedPassword := salt + password
	hash, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if string(hash) != hashStored {
		return ErrorWrongPassword
	}
	return nil
}
