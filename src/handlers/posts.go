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
		err = errors.Join(utils.GetFunctionName(), err)
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	data.(state.StateController).WriteResponse()
}
