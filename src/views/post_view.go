package views

import "forum/src/models"

func PostView(data models.ResponseStruct) {
	data.SetView("post_view").WriteResponse()
}
