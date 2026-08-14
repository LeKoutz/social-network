package handlers

import (
	"forum/src/controllers"
	"errors"
	"forum/src/parsers"
	"forum/src/state"
	"forum/src/utils"
)

func HandleShowCategory(data state.StateHandler) {
	var err error
	data.EditCategory().Id, err = parsers.ParseCategoryId(data)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		data.SetErrorConsume(err).WriteResponse()
		return
	}
	err = controllers.ShowCategory(data.(state.StateController))
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		data.SetErrorConsume(err).WriteResponse()
		return
	}
	data.WriteResponse()
}
