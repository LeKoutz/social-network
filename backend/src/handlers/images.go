package handlers

import (
	"errors"
	"forum/src/controllers"
	"forum/src/ferror"
	"forum/src/state"
	"forum/src/utils"
	"net/http"
)

func HandleImages(data state.StateHandler) {
	imgURL, err := controllers.HandleImages(data.(state.StateController))
	if err != nil {
		err = ferror.ErrorNotFound
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		data.SetErrorConsume(err).WriteResponse()
		return
	}
	http.ServeFile(*data.EditResponse(), data.GetRequest(), "./uploads/images/"+imgURL)
}
