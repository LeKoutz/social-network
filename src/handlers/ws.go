package handlers

import (
	"forum/src/controllers"
	"forum/src/state"
)

func HandleWs(data state.StateHandler) {
	err := controllers.ServeWs(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(err).WriteResponse()
		return
	}
}
