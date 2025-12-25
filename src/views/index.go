package views

import (
	"forum/src/models"
)

func Index(data models.ResponseStruct) {
	data.SetView("index_view").WriteResponse()
}
