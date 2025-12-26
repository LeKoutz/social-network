package controllers

import (
	"forum/src/models"
	"forum/src/views"
	"net/http"
	"regexp"
	"time"

	"github.com/gofrs/uuid"
)

var (
	AdminUser = models.User{
		Username: "admin",
		Role:     "admin",
		LoggedIn: true,
	}
)

func userLogin(data models.ResponseStruct) {
	if data.User.LoggedIn {
		Index(*data.SetErrorConsume(models.ErrorAlreadyLoggedIn))
		return
	}
	views.UserLogin(data)
}

func userLogout(data models.ResponseStruct) {
	if data.Request.Method != http.MethodGet {
		data.SetErrorConsume(models.ErrorMethodNotAllowed).WriteResponse()
		return
	}
	GuestUser := models.GetGuestUser()
	cookie, err := data.Request.Cookie("__Host-FRMSessionID")
	if err != nil {
		data.SetUser(GuestUser).SetErrorConsume(err)
		data.SetView("error_view").WriteResponse()
		return
	}
	user, err := models.GetUserBySession(cookie.Value)
	if err != nil {
		data.SetUser(user)
		data.SetErrorConsume(err)
		data.SetView("user_register_view").WriteResponse()
		return
	}
	err = models.SetUserSession(user.Id, "")
	if err != nil {
		data.SetUser(user)
		data.SetErrorConsume(err)
		data.SetView("user_register_view").WriteResponse()
		return
	}
	http.SetCookie(data.Response, nullifyCookie(cookie))
	data.SetUser(GuestUser)
	views.UserLogout(data)
}

func nullifyCookie(cookie *http.Cookie) *http.Cookie {
	cookie = &http.Cookie{
		Name:     "__Host-FRMSessionID",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}
	return cookie
}

func attemptLogin(data models.ResponseStruct) {
	var email string
	var password string
	var err error
	err = data.Request.ParseForm()
	if err != nil {
		data.User = models.GetGuestUser()
		data.SetErrorConsume(err)
		data.SetView("user_register_view").WriteResponse()
		return
	}
	if len(data.Request.Form.Get("email")) != 0 {
		email = data.Request.Form.Get("email")
	}
	if len(data.Request.Form.Get("password")) != 0 {
		password = data.Request.Form.Get("password")
	}
	err = Auth(email, password)
	if err != nil {
		data.User = models.GetGuestUser()
		data.SetErrorConsume(err)
		data.SetView("user_register_view").WriteResponse()
		return
	}
	sessionValue, err := uuid.NewV4()
	if err != nil {
		data.User = models.GetGuestUser()
		data.SetErrorConsume(err)
		data.SetView("user_register_view").WriteResponse()
		return
	}
	data.User, err = models.GetUserByEmail(email)
	if err != nil {
		data.User = models.GetGuestUser()
		data.SetErrorConsume(err)
		data.SetView("user_register_view").WriteResponse()
		return
	}
	data.User.LoggedIn = true
	err = models.SetUserSession(data.User.Id, sessionValue.String())
	if err != nil {
		data.User = models.GetGuestUser()
		data.SetErrorConsume(err)
		data.SetView("user_register_view").WriteResponse()
		return
	}
	cookie := &http.Cookie{
		Name:     "__Host-FRMSessionID",
		Value:    sessionValue.String(),
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSite(http.SameSiteStrictMode),
	}
	http.SetCookie(data.Response, cookie)
	data.SetView("index_view").WriteResponse()
}

func registerUser(data models.ResponseStruct) {
	if data.User.LoggedIn {
		Index(*data.SetErrorConsume(models.ErrorAlreadyLoggedIn))
		return
	}
	var err error
	// var user models.User
	data.User.Username = data.Request.FormValue("username")
	data.User.Email = data.Request.FormValue("email")
	if err = data.User.ValidateUser(); err != nil {
		data.SetUser(data.User).SetErrorConsume(err)
		data.SetView("user_register_view").WriteResponse()
		return
	}
	if models.IsEmailRegistered(data.User.Email) {
		data.SetUser(data.User).SetErrorConsume(models.ErrorEmailIsRegistered)
		data.SetView("user_register_view").WriteResponse()
		return
	}
	if !models.IsUniqueUsername(data.User.Username) {
		data.SetUser(data.User).SetErrorConsume(models.ErrorUsernameTaken)
		data.SetView("user_register_view").WriteResponse()
		return
	}
	if !CompareRegistrationPasswords(data.Request.FormValue("password1"), data.Request.FormValue("password2")) {
		data.SetUser(data.User).SetErrorConsume(models.ErrorPasswordMismatch)
		data.SetView("user_register_view").WriteResponse()
		return
	}
	password := data.Request.FormValue("password1")
	if err = validatePasswordStrength(password); err != nil {
		data.SetUser(data.User).SetErrorConsume(err)
		data.SetView("user_register_view").WriteResponse()
		return
	}
	data.User.Hash, err = HashPassword(password)
	if err != nil {
		data.SetUser(data.User).SetErrorConsume(err)
		data.SetView("user_register_view").WriteResponse()
		return
	}
	if err = models.RegisterUserOnDB(data.User); err != nil {
		data.SetUser(data.User).SetErrorConsume(models.ErrorInvalidUser)
		data.SetView("user_register_view").WriteResponse()
		return
	}
	// if err = Auth(user.Email, password); err != nil {
	// 	data.SetUser(user).SetError(*(&models.Error{}).Consume(err))
	// 	data.SetView("user_register_view").WriteResponse()
	// 	return
	// }
	// user.LoggedIn = true
	views.Index(data)
}

// Strong password validation. Makes sure the password is in between 10-16
// characters and includes letters, numbers and punctation symbols
func validatePasswordStrength(password string) error {
	unameMask := regexp.MustCompile(`^[[:punct:][:alnum:]]{10,16}$`)
	if !unameMask.MatchString(password) {
		return models.ErrorWeakPassword
	}
	return nil
}

func showUserPosts(data models.ResponseStruct) {
	posts, err := models.GetPostsByUserId(data.User.Id)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	for i := range posts {
		err = posts[i].GetReactions()
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
			return
		}
		err = posts[i].GetReactionsByUserId(data.User.Id)
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
			return
		}
	}
	data.Posts = posts
	data.View = "posts_view"
	data.WriteResponse()
}

func showUserLikedPosts(data models.ResponseStruct) {
	var err error
	var posts models.Posts
	posts, err = models.GetLikedPostsByUserId(data.User.Id)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	for i := range posts {
		err = posts[i].GetReactions()
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
			return
		}
		err = posts[i].GetReactionsByUserId(data.User.Id)
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
			return
		}
	}
	data.Posts = posts
	data.View = "posts_view"
	data.WriteResponse()
}
