package views

import (
	"forum/src/models"
)

func UserLogin(data models.ResponseStruct) {
	data.SetView("user_login_view").WriteResponse()
}
