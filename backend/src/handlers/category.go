package handlers

import (
	"forum/src/controllers"
	"forum/src/ferror"
	"forum/src/parsers"
	"forum/src/state"
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
