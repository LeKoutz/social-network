package controllers

import (
	"forum/src/models"
)

func Index(data models.ResponseStruct) {
	if data.User.LoggedIn {
		categories, err := models.GetAllCategories()
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
			return
		}
		data.SetCategories(categories)
	}
	data.WriteResponse()
}
