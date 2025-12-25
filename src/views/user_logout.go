package views

import (
	"forum/src/models"
	"net/http"
	"time"
)

func UserLogout(res http.ResponseWriter, req *http.Request, _ models.User) {
	GuestUser := models.GetGuestUser()
	data := models.ResponseStruct{}
	data.Init()
	cookie, err := req.Cookie("__Host-FRMSessionID")
	if err != nil {
		data.SetUser(GuestUser).SetError(*(&models.Error{}).Consume(err))
		data.SetView("error_view").WriteResponse(res)
		return
	}
	user, err := models.GetUserBySession(cookie.Value)
	if err != nil {
		data.SetUser(user)
		data.SetError(*(&models.Error{}).Consume(err))
		data.SetView("user_register_view").WriteResponse(res)
		return
	}
	err = models.SetUserSession(user.Id, "")
	if err != nil {
		data.SetUser(user)
		data.SetError(*(&models.Error{}).Consume(err))
		data.SetView("user_register_view").WriteResponse(res)
		return
	}
	http.SetCookie(res, nullifyCookie(cookie))
	data.SetUser(GuestUser)
	data.SetView("user_logout_view").WriteResponse(res)
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

