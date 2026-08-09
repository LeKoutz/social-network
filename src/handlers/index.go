package handlers

import (
	"forum/src/controllers"
	"forum/src/state"
)

func HandleIndex(data state.StateHandler) {
	err := controllers.Index(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	data.(state.StateController).WriteResponse()
}
