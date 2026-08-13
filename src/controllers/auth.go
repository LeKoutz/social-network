package controllers

import (
	"errors"
	"forum/src/ferror"
	"forum/src/models"
	"forum/src/utils"

	"golang.org/x/crypto/bcrypt"
)

func Auth(identifier, password string) error {
	var err error
	if !models.IsEmailRegistered(identifier) && !models.IsUsernameRegistered(identifier) {
		return ferror.ErrorNotRegistered
	}
	var user models.UserType
	err = user.SelectUserPasswordByIdentifier(identifier)
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ferror.ErrorWrongPassword
		}
		return errors.Join(utils.GetFunctionName(), err)
	}
	return nil
}
