package handlers

import (
	"forum/src/controllers"
	"forum/src/ferror"
	"forum/src/state"
)

func HandleShowGroups(data state.StateHandler) {
	err := controllers.ShowGroups(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.WriteResponse()
}
