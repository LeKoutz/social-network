package models

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

var (
	ErrorNotRegistered           = errors.New("Email is not registered")
	ErrorEmailIsRegistered       = errors.New("Email is already registered")
	ErrorInvalidUsername         = errors.New("Username is invalid")
	ErrorUsernameTaken           = errors.New("Username is taken")
	ErrorInvalidUser             = errors.New("Invalid user")
	ErrorWrongPassword           = errors.New("Wrong password")
	ErrorWeakPassword            = errors.New("Weak password. Use lower and upper case letters, symbols and number. Length must be between 10-16 characters.")
	ErrorPasswordMismatch        = errors.New("Password mismatch")
	ErrorNotFound                = errors.New("Not found")
	ErrorPostEmptyId             = errors.New("Post ID can't be empty.")
	ErrorInvalidPostId           = errors.New("Invalid post ID")
	ErrorInvalidCommentId        = errors.New("Invalid comment ID")
	ErrorInvalidCategoryId       = errors.New("Invalid category ID")
	ErrorPostBodyEmpty           = errors.New("Post body can't be empty.")
	ErrorPostTitleEmpty          = errors.New("Post title can't be empty.")
	ErrorPostHasNoCategory       = errors.New("Post category can't be empty.")
	ErrorPostPermissionDenied    = errors.New("You must be logged in to create a post.")
	ErrorCommentEmpty            = errors.New("Comment can't be empty.")
	ErrorCommentTooLong          = errors.New("Comment is too long.")
	ErrorCommentPermissionDenied = errors.New("You must be logged in to create a comment.")
	ErrorCommentEmptyId          = errors.New("Comment ID can't be empty.")
	ErrorCategoryEmptyId         = errors.New("Category ID can't be empty.")
	ErrorCategoryAlreadyExists   = errors.New("Category already exists.")
	ErrorCategoryNameEmpty       = errors.New("Category name can't be empty.")
	ErrorCategoryNameTooLong     = errors.New("Category name is too long. Use less than 128 characters.")
	ErrorUnauthorizedAction      = errors.New("Unauthorized action.")
	ErrorMethodNotAllowed        = errors.New("Method not allowed.")
	ErrorBadRequest              = errors.New("Bad request.")
)

type Error struct {
	Has        bool
	StatusCode int
	Message    string
	Error      error
}

type ErrorIface interface {
	LogError()
	RespondError(http.ResponseWriter)
	Consume(error) Error
}

// Converts error type to *Error
func (e *Error) Consume(err error) *Error {
	e.Message = strings.ReplaceAll(err.Error(), "\n", ": ")
	e.Error = err
	e.Has = true
	switch err {
	case ErrorNotFound:
		e.StatusCode = http.StatusNotFound
	case
		ErrorUnauthorizedAction,
		ErrorCommentPermissionDenied,
		ErrorPostPermissionDenied:
		e.StatusCode = http.StatusForbidden
	case ErrorMethodNotAllowed:
		e.StatusCode = http.StatusMethodNotAllowed
	}
	return e
}

// Logs *Error to terminal
func (e *Error) LogError() {
	log.Printf("Error: %s", e.Message)
}

// Responds *Error to user with the error_view template
func (e *Error) RespondError(res http.ResponseWriter, user User) {
	data := ResponseStruct{}
	data.Init().SetResponse(res)
	data.SetError(*e).SetUser(user).SetView("error_view").WriteResponse()
}

// Logs *Error to terminal and responds to user with error_view template
func (e *Error) LogAndRespondError(res http.ResponseWriter, user User) {
	e.LogError()
	e.RespondError(res, user)
}
