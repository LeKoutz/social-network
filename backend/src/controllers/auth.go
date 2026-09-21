package controllers

import (
	"errors"
	"forum/src/ferror"
	"forum/src/models"

	"golang.org/x/crypto/bcrypt"
)

func Auth(identifier, password string) error {
	var err error
	if !models.IsEmailRegistered(identifier) && !models.IsUsernameRegistered(identifier) {
		err = ferror.ErrorNotRegistered
		return ferror.ReturnErr(err)
	}
	var user models.UserType
	user.Identifier = identifier
	err = user.GetUserPasswordByIdentifier()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			err = ferror.ErrorWrongPassword
		}
		return ferror.ReturnErr(err)
	}
	return nil
}
