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
		err = ferror.ErrorNotFound
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	http.ServeFile(*data.EditResponse(), data.GetRequest(), "./uploads/images/"+imgURL)
}
