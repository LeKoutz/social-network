package forum

import "net/http"

var (
	GuestUser = User{
		Name:     "guest",
		Role:     "guest",
		LoggedIn: false,
	}
	AdminUser = User{
		Name:     "admin",
		Role:     "admin",
		LoggedIn: true,
	}
)

type User struct {
	Name          string
	Email         string
	Role          string
	LoggedIn      bool
	OwnedPosts    Posts
	OwnedComments Comments
	OwnedLikes    Likes
	OwnedDislikes Dislikes
}

func showLogin(res http.ResponseWriter, _ *http.Request, user User) {
	data := ValuesToClient()
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
		Name:  "access",
		Value: "admin",
		Path:  "/",
	}
	http.SetCookie(res, cookie)
	data := ValuesToClient()
	data.User = AdminUser
	data.User.Name = username
	respondView(res, "index_view", data)
}

func showRegister(res http.ResponseWriter, _ *http.Request, user User) {
	data := ValuesToClient()
	data.User = user
	respondView(res, "user_register_view", data)
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
	data := ValuesToClient()
	data.User = GuestUser
	respondView(res, "user_logout_view", data)
}
