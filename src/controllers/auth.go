package controllers

import (
	"errors"
	"forum/src/models"

	"golang.org/x/crypto/bcrypt"
)

func CompareRegistrationPasswords(pass1, pass2 string) bool {
	return pass1 == pass2
}

func Auth(email, password string) error {
	var err error
	// I guess in order to authenticate against an email and a password, we will
	// need to first check if email is registered
	if !models.IsEmailRegistered(email) {
		// Return an error... this should be sent back to umh... places...?!
		return models.ErrorNotRegistered
	}
	var user models.User
	user, err = models.GetUserByEmail(email)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return models.ErrorWrongPassword
		}
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
