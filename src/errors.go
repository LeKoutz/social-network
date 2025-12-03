package forum

import (
	"errors"
	"log"
	"net/http"
)

var (
	ErrorNotRegistered     = errors.New("Email is not registered")
	ErrorEmailIsRegistered = errors.New("Email is already registered")
	ErrorInvalidUsername   = errors.New("Username is invalid")
	ErrorUsernameTaken     = errors.New("Username is taken")
	ErrorInvalidUser       = errors.New("Invalid user")
	ErrorWrongPassword     = errors.New("Wrong password")
	ErrorWeakPassword      = errors.New("Weak password. Use lower and upper case letters, symbols and number. Length must be between 10-16 characters.")
	ErrorPasswordMismatch  = errors.New("Password mismatch")
	ErrorNotFound          = errors.New("Not found")
)

type Error struct {
	Has     bool
	Message string
	Error   error
}

type ErrorIface interface {
	LogError()
	RespondError(http.ResponseWriter)
	Consume(error) Error
}

func (e *Error) Consume(err error) Error {
	e.Message = err.Error()
	e.Error = err
	e.Has = true
	return *e
}

func (e *Error) LogError() {
	log.Printf("Error: %s", e.Message)
}

func (e *Error) RespondError(res http.ResponseWriter) {
	data := ReturnMockResponse()
	data.Error = *e
	respondView(res, "error_view", data)
}
