package forum

import (
	"errors"
	"log"
	"net/http"
)

var (
	ErrorNotRegistered = errors.New("Email is not registered")
	ErrorWrongPassword = errors.New("Wrong password")
	ErrorNotFound      = errors.New("Not found")
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
	data := ValuesToClient()
	data.Error = *e
	respondView(res, "error_view", data)
}
