package views

import (
	"forum/src/models"
)

func UserRegister(data models.ResponseStruct) {
	data.SetView("user_register_view").WriteResponse()
}
