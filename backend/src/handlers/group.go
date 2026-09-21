package handlers

import (
	"forum/src/controllers"
	"forum/src/ferror"
	"forum/src/parsers"
	"forum/src/state"
	"net/http"
)

func HandleCreateGroup(data state.StateHandler) {
	if data.GetRequest().Method != http.MethodPost {
		data.SetErrorConsume(ferror.ReturnErr(ferror.ErrorMethodNotAllowed))
		data.WriteResponse()
		return
	}
	var err error
	err = parsers.ParseCreateGroupRequest(data)
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	err = controllers.CreateGroup(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.WriteResponse()
}

func HandleShowGroup(data state.StateHandler) {
	var err error
	data.EditGroup().Id, err = parsers.ParseGroupId(data)
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	err = controllers.ShowGroup(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.WriteResponse()
}
