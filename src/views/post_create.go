package views

import (
	"forum/src/models"
)

func PostCreate(data models.ResponseStruct) {
	categories, err := models.GetAllCategories()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	data.SetCategories(categories)
	data.SetView("post_create_view").WriteResponse()
}
