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

func AttemptRegister(data state.StateController) {
	var err error
	if err = data.EditUser().ValidateUser(); err != nil {
		data.SetUser(data.GetUser()).SetErrorConsume(err)
		data.WriteResponse()
		return
	}
	if models.IsEmailRegistered(data.GetUser().Email) {
		data.SetUser(data.GetUser()).SetErrorConsume(ferror.ErrorEmailIsRegistered)
		data.WriteResponse()
		return
	}
	if !models.IsUniqueUsername(data.GetUser().Username) {
		data.SetUser(data.GetUser()).SetErrorConsume(ferror.ErrorUsernameTaken)
		data.WriteResponse()
		return
	}
	if !CompareRegistrationPasswords(data.GetRequest().FormValue("password1"), data.GetRequest().FormValue("password2")) {
		data.SetUser(data.GetUser()).SetErrorConsume(ferror.ErrorPasswordMismatch)
		data.WriteResponse()
		return
	}
	password := data.GetRequest().FormValue("password1")
	if err = ValidatePasswordStrength(password); err != nil {
		data.SetUser(data.GetUser()).SetErrorConsume(err)
		data.WriteResponse()
		return
	}
	data.EditUser().Hash, err = utils.HashPassword(password)
	if err != nil {
		data.SetUser(data.GetUser()).SetErrorConsume(err)
		data.WriteResponse()
		return
	}
	data.EditUser().FirstName = data.GetRequest().FormValue("first_name")
	data.EditUser().LastName = data.GetRequest().FormValue("last_name")
	age, err := utils.StringToInt64(data.GetRequest().FormValue("age"))
	if err != nil {
		data.SetUser(data.GetUser()).SetErrorConsume(err)
		data.WriteResponse()
		return
	}
	data.EditUser().Age = age
	data.EditUser().Gender = data.GetRequest().FormValue("gender")
	if err = data.EditUser().Add(); err != nil {
		data.SetUser(data.GetUser()).SetErrorConsume(ferror.ErrorInvalidUser)
		data.WriteResponse()
		return
	}
	data.SetUser(models.GetGuestUser())
	data.SetMessage(models.Message{
		Content: "Registration was successful",
		Type:    "Success",
		Has:     true,
	})
	data.WriteResponse()
}

func AttemptLogin(data state.StateController) error {
	var identifier string
	var password string
	var err error
	if len(data.GetRequest().Form.Get("identifier")) != 0 {
		identifier = data.GetRequest().Form.Get("identifier")
	} else {
		return ferror.ErrorEmailFieldEmpty
	}
	if len(data.GetRequest().Form.Get("password")) != 0 {
		password = data.GetRequest().Form.Get("password")
	} else {
		return ferror.ErrorPasswordFieldEmpty
	}
	err = Auth(identifier, password)
	if err != nil {
		data.SetUser(models.GetGuestUser())
		if !errors.Is(err, ferror.ErrorWrongPassword) && !errors.Is(err, ferror.ErrorNotRegistered) {
			return ferror.ErrorInternalServerError
		}
		return errors.Join(utils.GetFunctionName(), err)
	}
	sessionValue, err := uuid.NewV4()
	if err != nil {
		data.SetUser(models.GetGuestUser())
		return errors.Join(utils.GetFunctionName(), err)
	}
	data.EditUser().SessionId = sessionValue.String()
	data.EditUser().Identifier = identifier
	err = data.EditUser().GetUserByIdentifier(data.EditUser().Identifier)
	if err != nil {
		data.SetUser(models.GetGuestUser())
		return errors.Join(utils.GetFunctionName(), err)
	}
	data.EditUser().LoggedIn = true
	err = data.EditUser().SetUserSession(sessionValue.String())
	if err != nil {
		data.SetUser(models.GetGuestUser())
		data.SetErrorConsume(err)
		data.WriteResponse()
		return err
	}
	cookie := &http.Cookie{
		Name:     "__Host-FRMSessionID",
		Value:    sessionValue.String(),
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSite(http.SameSiteStrictMode),
	}
	http.SetCookie(*data.EditResponse(), cookie)
	data.SetMessage(
		models.Message{
			Has:     true,
			Type:    "Success",
			Content: "Login successful",
		},
	)
	return nil
}

// Strong password validation. Makes sure the password is in between 10-16
// characters and includes letters, numbers and/or punctuation symbols
func ValidatePasswordStrength(password string) error {
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
