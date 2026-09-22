package handlers

import (
	"forum/src/controllers"
	"forum/src/parsers"
	"forum/src/state"
	"forum/src/ferror"
)


func HandleShowProfileView(data state.StateHandler) {
	var err error
	data.EditProfile().UserId, err = parsers.ParseProfileId(data)
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	err = controllers.ShowProfile(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.WriteResponse()
}