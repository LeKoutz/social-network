package forum

import (
	"net/http"
	"net/mail"
	"regexp"
)

var (
	GuestUser = User{
		Username: "guest",
		Role:     "guest",
		LoggedIn: false,
	}
	AdminUser = User{
		Username: "admin",
		Role:     "admin",
		LoggedIn: true,
	}
)

type User struct {
	Username      string
	Salt          string
	Hash          string
	Email         string
	Role          string
	LoggedIn      bool
	OwnedPosts    Posts
	OwnedComments Comments
	OwnedLikes    Likes
	OwnedDislikes Dislikes
}

func showLogin(res http.ResponseWriter, _ *http.Request, user User) {
	data := ReturnMockResponse()
	data.User = user
	respondView(res, "user_login_view", data)
}

func GetUserSalt(email string) string {
	// select from users where email = email
	// row->salt
	return ""
}

func GetUserHash(email string) string {
	// select from users where email = email
	// row->hash
	return ""
}

func attemptLogin(res http.ResponseWriter, req *http.Request, _ User) {
	var username string
	if len(req.Form.Get("email")) != 0 {
		username = req.Form.Get("email")
	}
	cookie := &http.Cookie{
		Name:     "access",
		Value:    "admin",
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(res, cookie)
	data := ReturnMockResponse()
	data.User = AdminUser
	data.User.Username = username
	respondView(res, "index_view", data)
}

func showRegister(res http.ResponseWriter, _ *http.Request, user User) {
	data := ReturnMockResponse()
	data.User = user
	respondView(res, "user_register_view", data)
}

func registerUser(res http.ResponseWriter, req *http.Request) {
	var err error
	var user User
	user.Username = req.FormValue("username")
	user.Email = req.FormValue("email")
	if err = user.validateUser(); err != nil {
		data := ReturnMockResponse()
		var e Error
		data.User = user
		data.Error = e.Consume(err)
		respondView(res, "user_register_view", data)
	}
	if IsEmailRegistered(user.Email) {
		data := ReturnMockResponse()
		var e Error
		data.User = user
		data.Error = e.Consume(ErrorEmailIsRegistered)
		respondView(res, "user_register_view", data)
		return
	}
	if !IsUniqueUsername(user.Username) {
		data := ReturnMockResponse()
		var e Error
		data.User = user
		data.Error = e.Consume(ErrorUsernameTaken)
		respondView(res, "user_register_view", data)
		return
	}
	if !CompareRegistrationPasswords(req.FormValue("password1"), req.FormValue("password2")) {
		data := ReturnMockResponse()
		var e Error
		data.User = user
		data.Error = e.Consume(ErrorPasswordMismatch)
		respondView(res, "user_register_view", data)
		return
	}
	password := req.FormValue("password1")
	if err = validateStrongPassword(password); err != nil {
		data := ReturnMockResponse()
		var e Error
		data.User = user
		data.Error = e.Consume(err)
		respondView(res, "user_register_view", data)
		return
	}
	user.Salt = SaltGenerator(26 - len([]byte(password)))
	user.Hash, err = HashPassword(SaltPassword(user.Salt, password))
	if err != nil {
		data := ReturnMockResponse()
		var e Error
		data.User = user
		data.Error = e.Consume(err)
		respondView(res, "user_register_view", data)
		return
	}
	if err = registerUserOnDB(user); err != nil {
		data := ReturnMockResponse()
		var e Error
		data.User = user
		data.Error = e.Consume(ErrorInvalidUser)
		respondView(res, "user_register_view", data)
		return
	}
	showIndex(res, req, user)
}

func nullifyCookie(cookie *http.Cookie) {
	cookie.Value = ""
}

func showLogout(res http.ResponseWriter, req *http.Request, _ User) {
	cookies := req.Cookies()
	for i := range cookies {
		nullifyCookie(cookies[i])
	}
	if len(cookies) != 0 {
		http.SetCookie(res, cookies[0])
	} else {
		respondView(res, "error_view", ResponseStruct{
			Error: Error{
				Has:     true,
				Message: "Lol",
			},
		})
		return
	}
	data := ReturnMockResponse()
	data.User = GuestUser
	respondView(res, "user_logout_view", data)
}

func (u *User) validateUsername() error {
	unameMask := regexp.MustCompile(`^[a-zA-Z0-9]{4,15}$`)
	if !unameMask.MatchString((*u).Username) {
		return ErrorInvalidUsername
	}
	return nil
}

func (u *User) validateEmail() error {
	_, err := mail.ParseAddress(u.Email)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) validateUser() error {
	var err error
	if err = u.validateUsername(); err != nil {
		return err
	}
	if err = u.validateEmail(); err != nil {
		return err
	}
	return nil
}

func validateStrongPassword(password string) error {
	unameMask := regexp.MustCompile(`^[[:punct:][:alnum:]]{10,16}$`)
	if !unameMask.MatchString(password) {
		return ErrorWeakPassword
	}
	return nil
}
