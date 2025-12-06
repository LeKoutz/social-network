package forum

import (
	"errors"
	"log"
	"net/http"
)

var (
	ErrorNotRegistered         = errors.New("Email is not registered")
	ErrorEmailIsRegistered     = errors.New("Email is already registered")
	ErrorInvalidUsername       = errors.New("Username is invalid")
	ErrorUsernameTaken         = errors.New("Username is taken")
	ErrorInvalidUser           = errors.New("Invalid user")
	ErrorWrongPassword         = errors.New("Wrong password")
	ErrorWeakPassword          = errors.New("Weak password. Use lower and upper case letters, symbols and number. Length must be between 10-16 characters.")
	ErrorPasswordMismatch      = errors.New("Password mismatch")
	ErrorNotFound              = errors.New("Not found")
	ErrorPostEmptyId           = errors.New("Post ID can't be empty.")
	ErrorPostBodyEmpty         = errors.New("Post body can't be empty.")
	ErrorPostTitleEmpty        = errors.New("Post title can't be empty.")
	ErrorPostHasNoCategory     = errors.New("Post category can't be empty.")
	ErrorPostPermissionDenied  = errors.New("You must be logged in to create a post.")
	ErrorCategoryEmptyId       = errors.New("Category ID can't be empty.")
	ErrorCategoryAlreadyExists = errors.New("Category already exists.")
	ErrorCategoryNameEmpty     = errors.New("Category name can't be empty.")
	ErrorCategoryNameTooLong   = errors.New("Category name is too long. Use less than 128 characters.")
	ErrorUnauthorizedAction    = errors.New("Unauthorized action.")
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

func (e *Error) Consume(err error) *Error {
	e.Message = err.Error()
	e.Error = err
	e.Has = true
	return e
}

func (e *Error) LogError() {
	log.Printf("Error: %s", e.Message)
}

func (e *Error) RespondError(res http.ResponseWriter, user User) {
	data := ReturnMockResponse()
	data.Error = *e
	data.User = user
	respondView(res, "error_view", data)
}

func (e *Error) LogAndRespondError(res http.ResponseWriter, user User) {
	e.LogError()
	e.RespondError(res, user)
}
