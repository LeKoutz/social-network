package handlers

import (
	"forum/src/controllers"
	"forum/src/state"
	"forum/src/ferror"
)

func HandleShowPosts(data state.StateHandler) {
	err := controllers.ShowPosts(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.WriteResponse()
}
