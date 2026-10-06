package handlers

import (
	"forum/src/controllers"
	"forum/src/ferror"
	"forum/src/state"
)

func HandleIndex(data state.StateHandler) {
	err := controllers.Index(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.WriteResponse()
}
