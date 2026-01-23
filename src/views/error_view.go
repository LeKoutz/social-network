package views

import (
	"forum/src/models"
)

func Error(data models.ResponseStruct4ViewsIface) {
	data.SetView("error_view").WriteResponse()
}
