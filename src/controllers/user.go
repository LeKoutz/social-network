package controllers

import (
	"forum/src/models"
	"forum/src/views"
	"net/http"
	"regexp"

	"github.com/gofrs/uuid"
)

var (
	AdminUser = models.User{
		Username: "admin",
		Role:     "admin",
		LoggedIn: true,
	}
)

func attemptLogin(res http.ResponseWriter, req *http.Request, _ models.User) {
	var email string
	var password string
	var err error
	data := models.ResponseStruct{}
	data.Init()
	err = req.ParseForm()
	if err != nil {
		data.User = models.GetGuestUser()
		data.Error = *(&models.Error{}).Consume(err)
		data.SetView("user_register_view").WriteResponse()
		return
	}
	if len(req.Form.Get("email")) != 0 {
		email = req.Form.Get("email")
	}
	if len(req.Form.Get("password")) != 0 {
		password = req.Form.Get("password")
	}
	err = Auth(email, password)
	if err != nil {
		data.User = models.GetGuestUser()
		data.Error = *(&models.Error{}).Consume(err)
		data.SetView("user_register_view").WriteResponse()
		return
	}
	sessionValue, err := uuid.NewV4()
	if err != nil {
		data.User = models.GetGuestUser()
		data.Error = *(&models.Error{}).Consume(err)
		data.SetView("user_register_view").WriteResponse()
		return
	}
	data.User, err = models.GetUserByEmail(email)
	if err != nil {
		data.User = models.GetGuestUser()
		data.Error = *(&models.Error{}).Consume(err)
		data.SetView("user_register_view").WriteResponse()
		return
	}
	data.User.LoggedIn = true
	err = models.SetUserSession(data.User.Id, sessionValue.String())
	if err != nil {
		data.User = models.GetGuestUser()
		data.Error = *(&models.Error{}).Consume(err)
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
	http.SetCookie(res, cookie)
	data.SetView("index_view").WriteResponse()
}

func registerUser(res http.ResponseWriter, req *http.Request, user models.User) {
	data := models.ResponseStruct{}
	data.Init().SetResponse(res)
	if user.LoggedIn {
		(&models.Error{}).Consume(models.ErrorBadRequest).LogAndRespondError(res, user)
	}
	var err error
	// var user models.User
	user.Username = req.FormValue("username")
	user.Email = req.FormValue("email")
	if err = user.ValidateUser(); err != nil {
		data.SetUser(user).SetError(*(&models.Error{}).Consume(err))
		data.SetView("user_register_view").WriteResponse()
		return
	}
	if models.IsEmailRegistered(user.Email) {
		data.SetUser(user).SetError(*(&models.Error{}).Consume(models.ErrorEmailIsRegistered))
		data.SetView("user_register_view").WriteResponse()
		return
	}
	if !models.IsUniqueUsername(user.Username) {
		data.SetUser(user).SetError(*(&models.Error{}).Consume(models.ErrorUsernameTaken))
		data.SetView("user_register_view").WriteResponse()
		return
	}
	if !CompareRegistrationPasswords(req.FormValue("password1"), req.FormValue("password2")) {
		data.SetUser(user).SetError(*(&models.Error{}).Consume(models.ErrorPasswordMismatch))
		data.SetView("user_register_view").WriteResponse()
		return
	}
	password := req.FormValue("password1")
	if err = validatePasswordStrength(password); err != nil {
		data.SetUser(user).SetError(*(&models.Error{}).Consume(err))
		data.SetView("user_register_view").WriteResponse()
		return
	}
	user.Hash, err = HashPassword(password)
	if err != nil {
		data.SetUser(user).SetError(*(&models.Error{}).Consume(err))
		data.SetView("user_register_view").WriteResponse()
		return
	}
	if err = models.RegisterUserOnDB(user); err != nil {
		data.SetUser(user).SetError(*(&models.Error{}).Consume(models.ErrorInvalidUser))
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

func showUserPosts(res http.ResponseWriter, user models.User) {
	data := models.ResponseStruct{}
	data.Init().SetResponse(res)
	data.User = user
	data.View = "posts_view"
	posts, err := models.GetPostsByUserId(user.Id)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	for i := range posts {
		err = posts[i].GetReactions()
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(res, user)
			return
		}
		err = posts[i].GetReactionsByUserId(user.Id)
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(res, user)
			return
		}
	}
	data.Posts = posts
	data.WriteResponse()
}

func showUserLikedPosts(res http.ResponseWriter, user models.User) {
	data := models.ResponseStruct{}
	data.Init().SetResponse(res)
	data.User = user
	data.View = "posts_view"
	var err error
	var posts models.Posts
	posts, err = models.GetLikedPostsByUserId(user.Id)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	for i := range posts {
		err = posts[i].GetReactions()
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(res, user)
			return
		}
		err = posts[i].GetReactionsByUserId(user.Id)
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(res, user)
			return
		}
	}
	data.Posts = posts
	data.WriteResponse()
}
