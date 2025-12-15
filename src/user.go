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
	Id            int64
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
	data := ResponseStruct{}
	data.Init().SetUser(user).SetView("user_login_view").WriteResponse(res)
}

func attemptLogin(res http.ResponseWriter, req *http.Request, _ User) {
	var email string
	var password string
	var err error
	data := ResponseStruct{}
	data.Init()
	err = req.ParseForm()
	if err != nil {
		data.User = GuestUser
		data.Error = *(&Error{}).Consume(err)
		data.SetView("user_register_view").WriteResponse(res)
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
		data.User = GuestUser
		data.Error = *(&Error{}).Consume(err)
		data.SetView("user_register_view").WriteResponse(res)
		return
	}
	sessionValue, err := uuid.NewV4()
	if err != nil {
		data.User = GuestUser
		data.Error = *(&Error{}).Consume(err)
		data.SetView("user_register_view").WriteResponse(res)
		return
	}
	data.User, err = getUserByEmail(email)
	if err != nil {
		data.User = GuestUser
		data.Error = *(&Error{}).Consume(err)
		data.SetView("user_register_view").WriteResponse(res)
		return
	}
	data.User.LoggedIn = true
	err = setUserSession(data.User.Id, sessionValue.String())
	if err != nil {
		data.User = GuestUser
		data.Error = *(&Error{}).Consume(err)
		data.SetView("user_register_view").WriteResponse(res)
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
	data.SetView("index_view").WriteResponse(res)
}

func showRegister(res http.ResponseWriter, _ *http.Request, user User) {
	data := ResponseStruct{}
	data.Init().SetUser(user)
	data.SetView("user_register_view").WriteResponse(res)
}

func registerUser(res http.ResponseWriter, req *http.Request) {
	var err error
	var user User
	data := ResponseStruct{}
	data.Init()
	user.Username = req.FormValue("username")
	user.Email = req.FormValue("email")
	if err = user.validateUser(); err != nil {
		data.SetUser(user).SetError(*(&Error{}).Consume(err))
		data.SetView("user_register_view").WriteResponse(res)
		return
	}
	if IsEmailRegistered(user.Email) {
		data.SetUser(user).SetError(*(&Error{}).Consume(ErrorEmailIsRegistered))
		data.SetView("user_register_view").WriteResponse(res)
		return
	}
	if !IsUniqueUsername(user.Username) {
		data.SetUser(user).SetError(*(&Error{}).Consume(ErrorUsernameTaken))
		data.SetView("user_register_view").WriteResponse(res)
		return
	}
	if !CompareRegistrationPasswords(req.FormValue("password1"), req.FormValue("password2")) {
		data.SetUser(user).SetError(*(&Error{}).Consume(ErrorPasswordMismatch))
		data.SetView("user_register_view").WriteResponse(res)
		return
	}
	password := req.FormValue("password1")
	if err = validatePasswordStrength(password); err != nil {
		data.SetUser(user).SetError(*(&Error{}).Consume(err))
		data.SetView("user_register_view").WriteResponse(res)
		return
	}
	user.Hash, err = HashPassword(password)
	if err != nil {
		data.SetUser(user).SetError(*(&Error{}).Consume(err))
		data.SetView("user_register_view").WriteResponse(res)
		return
	}
	if err = registerUserOnDB(user); err != nil {
		data.SetUser(user).SetError(*(&Error{}).Consume(ErrorInvalidUser))
		data.SetView("user_register_view").WriteResponse(res)
		return
	}
	if err = Auth(user.Email, password); err != nil {
		data.SetUser(user).SetError(*(&Error{}).Consume(err))
		data.SetView("user_register_view").WriteResponse(res)
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
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}
	return cookie
}

func showLogout(res http.ResponseWriter, req *http.Request, _ User) {
	data := ResponseStruct{}
	data.Init()
	cookie, err := req.Cookie("__Host-FRMSessionID")
	if err != nil {
		data.SetUser(GuestUser).SetError(*(&Error{}).Consume(err))
		data.SetView("error_view").WriteResponse(res)
		return
	}
	user, err := getUserBySession(cookie.Value)
	if err != nil {
		data.SetUser(user)
		data.SetError(*(&Error{}).Consume(err))
		data.SetView("user_register_view").WriteResponse(res)
		return
	}
	err = setUserSession(user.Id, "")
	if err != nil {
		data.SetUser(user)
		data.SetError(*(&Error{}).Consume(err))
		data.SetView("user_register_view").WriteResponse(res)
		return
	}
	http.SetCookie(res, nullifyCookie(cookie))
	data.SetUser(GuestUser)
	data.SetView("user_logout_view").WriteResponse(res)
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
func validatePasswordStrength(password string) error {
	unameMask := regexp.MustCompile(`^[[:punct:][:alnum:]]{10,16}$`)
	if !unameMask.MatchString(password) {
		return ErrorWeakPassword
	}
	return nil
}

// Check if user already liked this post
func hasUserAlreadyLikedPost(userId, postId int64) (bool, error) {
	existingReactionId, err := checkIfUserAlreadyLikedPost(userId, postId)
	if err != nil {
		return false, err
	}
	return existingReactionId != 0, nil
}

// Check if user already disliked this post
func hasUserAlreadyDislikedPost(userId, postId int64) (bool, error) {
	existingReactionId, err := checkIfUserAlreadyDislikedPost(userId, postId)
	if err != nil {
		return false, err
	}
	return existingReactionId != 0, nil
}

// Check if user already liked this comment
func hasUserAlreadyLikedComment(userId, commentId int64) (bool, error) {
	existingReactionId, err := checkIfUserAlreadyLikedPost(userId, commentId)
	if err != nil {
		return false, err
	}
	return existingReactionId != 0, nil
}

// Check if user already disliked this comment
func hasUserAlreadyDislikedComment(userId, commentId int64) (bool, error) {
	existingReactionId, err := checkIfUserAlreadyDislikedComment(userId, commentId)
	if err != nil {
		return false, err
	}
	return existingReactionId != 0, nil
}

func showUserPosts(res http.ResponseWriter, user User) {
	data := ResponseStruct{}
	data.Init()
	data.User = user
	data.View = "posts_view"
	posts, err := getPostsByUserId(user.Id)
	if err != nil {
		(&Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	for i := range posts {
		err = posts[i].getReactions()
		if err != nil {
			(&Error{}).Consume(err).LogAndRespondError(res, user)
			return
		}
		err = posts[i].getReactionsByUserId(user.Id)
		if err != nil {
			(&Error{}).Consume(err).LogAndRespondError(res, user)
			return
		}
	}
	data.Posts = posts
	data.WriteResponse(res)
}

func showUserLikedPosts(res http.ResponseWriter, user User) {
	data := ResponseStruct{}
	data.Init()
	data.User = user
	data.View = "posts_view"
	var err error
	var posts Posts
	posts, err = getLikedPostsByUserId(user.Id)
	if err != nil {
		(&Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	for i := range posts {
		err = posts[i].getReactions()
		if err != nil {
			(&Error{}).Consume(err).LogAndRespondError(res, user)
			return
		}
		err = posts[i].getReactionsByUserId(user.Id)
		if err != nil {
			(&Error{}).Consume(err).LogAndRespondError(res, user)
			return
		}
	}
	data.Posts = posts
	data.WriteResponse(res)
}
