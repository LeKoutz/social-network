package forum

import "errors"

var (
	ErrorNotRegistered = errors.New("Email is not registered")
	ErrorWrongPassword = errors.New("Wrong password")
)
