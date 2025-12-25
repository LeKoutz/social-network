package views

import "forum/src/models"

func PostsView(data models.ResponseStruct) {
	data.SetView("posts_view").WriteResponse()
}
