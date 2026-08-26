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
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		data.SetErrorConsume(err).WriteResponse()
		return
	}
	data.WriteResponse()
}
