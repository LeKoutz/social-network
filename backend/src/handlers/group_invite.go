package handlers

import (
	"forum/src/controllers"
	"forum/src/ferror"
	"forum/src/parsers"
	"forum/src/state"
)

func HandleInviteToGroup(data state.StateHandler) {
	var err error
	err = parsers.ParseInviteToGroupRequest(data)
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	err = controllers.InviteUserToGroup(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.WriteResponse()
}

func HandleJoinGroup(data state.StateHandler) {
	var err error
	err = parsers.ParseJoinGroupRequest(data)
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	err = controllers.RequestJoinGroup(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.WriteResponse()
}
