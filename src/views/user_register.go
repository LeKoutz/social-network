package views

import (
	"forum/src/models"
	"net/http"
)

func UserRegister(res http.ResponseWriter, _ *http.Request, user models.User) {
	data := models.ResponseStruct{}
	data.Init().SetUser(user).SetResponse(res)
	data.SetView("user_register_view").WriteResponse()
}
