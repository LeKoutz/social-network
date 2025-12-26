package views

import (
	"forum/src/models"
)

func UserLogout(data models.ResponseStruct) {
	data.SetView("user_logout_view").WriteResponse()
}
