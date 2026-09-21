package handlers

import (
	"forum/src/controllers"
	"forum/src/state"
	"forum/src/ferror"
)

func HandleIndex(data state.StateHandler) {
	err := controllers.Index(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.WriteResponse()
}
