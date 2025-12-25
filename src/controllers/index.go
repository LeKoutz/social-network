package controllers

import (
	"forum/src/models"
	"forum/src/views"
	"net/http"
)

func Index(res http.ResponseWriter, req *http.Request, user models.User) {
	data := models.ResponseStruct{}
	data.Init().SetUser(user).SetResponse(res).SetRequest(req)
	categories, err := models.GetAllCategories()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	data.SetCategories(categories)
	views.Index(data)
}
