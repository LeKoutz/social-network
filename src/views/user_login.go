package views

import (
	"forum/src/models"
	"net/http"
)

func UserLogin(res http.ResponseWriter, _ *http.Request, user models.User) {
	data := models.ResponseStruct{}
	data.Init().SetUser(user).SetView("user_login_view").WriteResponse(res)
}
