package views

import "forum/src/models"

func Categories(data models.ResponseStruct) {
	data.SetView("categories_view").WriteResponse()
}
