package forum

import (
	"net/http"
	"net/mail"
	"regexp"
	"time"

	"github.com/gofrs/uuid"
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
	Id            int
	Username      string
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

func GetUserHash(email string) string {
	// select from users where email = email
	// row->hash
	return ""
}

func attemptLogin(res http.ResponseWriter, req *http.Request, _ User) {
	var email string
	var password string
	var err error
	data := ReturnMockResponse()
	err = req.ParseForm()
	if err != nil {
		var e Error
		data.User = GuestUser
		data.Error = e.Consume(err)
		respondView(res, "user_register_view", data)
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
		var e Error
		data.User = GuestUser
		data.Error = e.Consume(err)
		respondView(res, "user_register_view", data)
		return
	}
	sessionValue, err := uuid.NewV4()
	if err != nil {
		var e Error
		data.User = GuestUser
		data.Error = e.Consume(err)
		respondView(res, "user_register_view", data)
		return
	}
	data.User, err = getUserByEmail(email)
	if err != nil {
		var e Error
		data.User = GuestUser
		data.Error = e.Consume(err)
		respondView(res, "user_register_view", data)
		return
	}
	data.User.LoggedIn = true
	err = setUserSession(data.User.Id, sessionValue.String())
	if err != nil {
		var e Error
		data.User = GuestUser
		data.Error = e.Consume(err)
		respondView(res, "user_register_view", data)
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
	user.Hash, err = HashPassword(password)
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
	if err = Auth(user.Email, password); err != nil {
		data := ReturnMockResponse()
		var e Error
		data.User = user
		data.Error = e.Consume(err)
		respondView(res, "user_register_view", data)
		return
	}
	user.LoggedIn = true
	showIndex(res, req, user)
}

func nullifyCookie(cookie *http.Cookie) *http.Cookie {
	cookie = &http.Cookie{
		Name:     "__Host-FRMSessionID",
		Value:    "",
		Path:     "/",
		Expires: time.Unix(0, 0),
		HttpOnly: true,
	}
	return cookie
}

func showLogout(res http.ResponseWriter, req *http.Request, _ User) {
	cookie, err := req.Cookie("__Host-FRMSessionID")
	if err != nil {
		respondView(res, "error_view", ResponseStruct{
			Error: Error{
				Has:     true,
				Message: "Lol",
			},
		})
		return
	}
	user, err := getUserBySession(cookie.Value)
	if err != nil {
		data := ReturnMockResponse()
		var e Error
		data.User = user
		data.Error = e.Consume(err)
		respondView(res, "user_register_view", data)
		return
	}
	err = setUserSession(user.Id, "")
	if err != nil {
		data := ReturnMockResponse()
		var e Error
		data.User = user
		data.Error = e.Consume(err)
		respondView(res, "user_register_view", data)
		return
	}
	http.SetCookie(res, nullifyCookie(cookie))
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

// Strong password validation. Makes sure the password is in between 10-16
// characters and includes letters, numbers and punctation symbols
func validateStrongPassword(password string) error {
	unameMask := regexp.MustCompile(`^[[:punct:][:alnum:]]{10,16}$`)
	if !unameMask.MatchString(password) {
		return ErrorWeakPassword
	}
	return nil
}

// Check if user already liked this post
func hasUserAlreadyLikedPost(userId, postId int) (bool, error) {
	existingReactionId, err := checkIfUserAlreadyLikedPost(userId, postId)
	if err != nil {
		return false, err
	}

	return existingReactionId != 0, nil
}

// Check if user already disliked this post
func hasUserAlreadyDislikedPost(userId, postId int) (bool, error) {
	existingReactionId, err := checkIfUserAlreadyDislikedPost(userId, postId)
	if err != nil {
		return false, err
	}

	return existingReactionId != 0, nil
}

// Check if user already liked this comment
func hasUserAlreadyLikedComment(userId, commentId int) (bool, error) {
	existingReactionId, err := checkIfUserAlreadyLikedPost(userId, commentId)
	if err != nil {
		return false, err
	}

	return existingReactionId != 0, nil
}

// Check if user already disliked this comment
func hasUserAlreadyDislikedComment(userId, commentId int) (bool, error) {
	existingReactionId, err := checkIfUserAlreadyDislikedComment(userId, commentId)
	if err != nil {
		return false, err
	}

	return existingReactionId != 0, nil
}
