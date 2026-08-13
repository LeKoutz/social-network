package controllers

import (
	"errors"
	"forum/src/ferror"
	"forum/src/models"
	"forum/src/state"
	"forum/src/utils"
	"net/http"
	"regexp"
	"time"

	"github.com/gofrs/uuid"
)

func UserLogout(data state.StateController) error {
	var err error
	GuestUser := models.GetGuestUser()
	cookie, err := data.GetRequest().Cookie("__Host-FRMSessionID")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return ferror.ErrorAlreadyLoggedOut
		}
		return errors.Join(utils.GetFunctionName(), err)
	}
	data.EditUser().SessionId = cookie.Value
	err = data.EditUser().GetUserBySession()
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	err = data.EditUser().SetUserSession("")
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	http.SetCookie(*data.EditResponse(), UnsetCookie())
	data.SetUser(GuestUser)
	data.SetMessage(models.Message{Has: true, Type: "Success", Content: "Logout successful"})
	return nil
}

func AttemptRegister(data state.StateController) error {
	var err error
	if data.GetUser().LoggedIn {
		return ferror.ErrorAlreadyLoggedIn
	}
	if err = validatePasswordStrength(data.GetUser().Password); err != nil {
		return err
	}
	if data.EditUser().Hash, err = utils.HashPassword(data.GetUser().Password); err != nil {
		return err
	}
	data.EditUser().Password = ""
	if err = data.EditUser().Add(); err != nil {
		return err
	}
	data.SetUser(models.GetGuestUser())
	data.SetMessage(models.Message{
		Content: "Registration was successful",
		Type:    "Success",
		Has:     true,
	})
	return nil
}

func AttemptLogin(data state.StateController) error {
	if data.GetUser().LoggedIn {
		return ferror.ErrorAlreadyLoggedIn
	}
	if err := authenticateUser(data); err != nil {
		return err
	}
	sessionValue, err := createUserSession(data)
	if err != nil {
		return err
	}
	setSessionCookie(data, sessionValue)
	data.SetMessage(
		models.Message{
			Has:     true,
			Type:    "Success",
			Content: "Login successful",
		},
	)
	return nil
}

func authenticateUser(data state.StateController) error {
	err := Auth(data.EditUser().Identifier, data.EditUser().Password)
	if err != nil {
		data.SetUser(models.GetGuestUser())
		if !errors.Is(err, ferror.ErrorWrongPassword) && !errors.Is(err, ferror.ErrorNotRegistered) {
			return errors.Join(utils.GetFunctionName(), ferror.ErrorInternalServerError)
		}
		return errors.Join(utils.GetFunctionName(), err)
	}
	data.EditUser().Password = ""
	return nil
}

func createUserSession(data state.StateController) (string, error) {
	sessionValue, err := uuid.NewV4()
	if err != nil {
		data.SetUser(models.GetGuestUser())
		return "", errors.Join(utils.GetFunctionName(), err)
	}
	data.EditUser().SessionId = sessionValue.String()
	err = data.EditUser().GetUserByIdentifier()
	if err != nil {
		data.SetUser(models.GetGuestUser())
		return "", errors.Join(utils.GetFunctionName(), err)
	}
	data.EditUser().Identifier = ""
	data.EditUser().LoggedIn = true
	err = data.EditUser().SetUserSession(sessionValue.String())
	if err != nil {
		data.SetUser(models.GetGuestUser())
		data.SetErrorConsume(err)
		return "", err
	}
	return sessionValue.String(), nil
}

func setSessionCookie(data state.StateController, sessionValue string) {
	cookie := &http.Cookie{
		Name:     "__Host-FRMSessionID",
		Value:    sessionValue,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSite(http.SameSiteStrictMode),
	}
	http.SetCookie(*data.EditResponse(), cookie)
}

// Strong password validation. Makes sure the password is in between 10-16
// characters and includes letters, numbers and/or punctuation symbols
func validatePasswordStrength(password string) error {
	unameMask := regexp.MustCompile(`^[[:punct:][:alnum:]]{10,16}$`)
	if !unameMask.MatchString(password) {
		return ferror.ErrorWeakPassword
	}
	return nil
}

func GetUserPosts(data state.StateController) error {
	posts, err := data.EditUser().GetPosts()
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	for i := range posts {
		err = posts[i].GetReactions()
		if err != nil {
			return errors.Join(utils.GetFunctionName(), err)
		}
		err = posts[i].GetReactionsByUserId(data.GetUser().Id)
		if err != nil {
			return errors.Join(utils.GetFunctionName(), err)
		}
	}
	// data.Posts = posts
	// data.WriteResponse()
	data.SetPosts(posts)
	return nil
}

func GetUserLikedPosts(data state.StateController) error {
	var err error
	var posts models.PostsType
	posts, err = data.EditUser().GetLikedPosts()
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	for i := range posts {
		err = posts[i].GetReactions()
		if err != nil {
			return errors.Join(utils.GetFunctionName(), err)
		}
		err = posts[i].GetReactionsByUserId(data.GetUser().Id)
		if err != nil {
			return errors.Join(utils.GetFunctionName(), err)
		}
	}
	data.SetPosts(posts)
	return nil
}

func GetUserActivity(data state.StateController) error {
	return data.EditUser().GetActivity()
}
