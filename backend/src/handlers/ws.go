package handlers

import (
	"forum/src/controllers"
	"forum/src/ferror"
	"forum/src/state"
)

func HandleWs(data state.StateHandler) {
	err := controllers.ServeWs(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
}
