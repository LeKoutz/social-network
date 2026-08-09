package handlers

import (
	"forum/src/controllers"
	"forum/src/ferror"
	"forum/src/state"
	"net/http"
)

func HandleImages(data state.StateHandler) {
	imgURL, err := controllers.HandleImages(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ErrorNotFound)
		data.(state.StateController).WriteResponse()
		return
	}
	http.ServeFile(*data.EditResponse(), data.GetRequest(), "./uploads/images/"+imgURL)
}
