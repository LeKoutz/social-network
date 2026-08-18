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
		err = ferror.ErrorNotRegistered
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	var user models.UserType
	err = user.SelectUserPasswordByIdentifier(identifier)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			err = ferror.ErrorWrongPassword
		}
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}
