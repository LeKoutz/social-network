package handlers

import (
	"forum/src/controllers"
	"forum/src/state"
	"forum/src/ferror"
)

func HandleWs(data state.StateHandler) {
	err := controllers.ServeWs(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
}
