package views

import "forum/src/models"

func Category(data models.ResponseStruct) {
	data.SetView("category_view").WriteResponse()
}
