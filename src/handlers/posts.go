package handlers

import (
	"errors"
	"forum/src/controllers"
	"forum/src/state"
	"forum/src/utils"
)

func HandleShowPosts(data state.StateHandler) {
	err := controllers.ShowPosts(data.(state.StateController))
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		data.SetErrorConsume(err).WriteResponse()
		return
	}
	data.WriteResponse()
}
