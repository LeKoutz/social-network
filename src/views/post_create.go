package views

import (
	"net/http"
	"forum/src/models"
)

func PostCreate(res http.ResponseWriter, _ *http.Request, user models.User) {
	data := models.ResponseStruct{}
	data.Init().SetUser(user)
	categories, err := models.GetAllCategories()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	data.SetCategories(categories)
	data.SetView("post_create_view").WriteResponse()
}
