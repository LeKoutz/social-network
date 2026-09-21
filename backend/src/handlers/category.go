package handlers

import (
	"forum/src/controllers"
	"forum/src/parsers"
	"forum/src/state"
	"forum/src/ferror"
)

func HandleShowCategory(data state.StateHandler) {
	var err error
	data.EditCategory().Id, err = parsers.ParseCategoryId(data)
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	err = controllers.ShowCategory(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.WriteResponse()
}
