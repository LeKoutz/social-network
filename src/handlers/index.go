package handlers

import (
	"errors"
	"forum/src/controllers"
	"forum/src/state"
	"forum/src/utils"
)

func HandleIndex(data state.StateHandler) {
	err := controllers.Index(data.(state.StateController))
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	data.(state.StateController).WriteResponse()
}
