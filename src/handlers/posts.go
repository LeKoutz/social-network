package handlers

import (
	"forum/src/controllers"
	"forum/src/state"
)

func HandleShowPosts(data state.StateHandler) {
	err := controllers.ShowPosts(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	data.(state.StateController).WriteResponse()
}
