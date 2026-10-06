package handlers

import (
	"forum/src/controllers"
	"forum/src/ferror"
	"forum/src/parsers"
	"forum/src/state"
)

func HandleShowFollows(data state.StateHandler) {
	err := controllers.GetFollowsForUser(data.(state.StateController))
	if err != nil {
		err = ferror.ErrorInternalServerError
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.WriteResponse()
}

func HandleAcceptFollowRequest(data state.StateHandler) {
	var err error
	err = parsers.ParseAcceptFollowRequest(data)
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	err = controllers.AcceptFollowRequest(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.WriteResponse()
}

func HandleDenyFollowRequest(data state.StateHandler) {
	var err error
	err = parsers.ParseDenyFollowRequest(data)
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	err = controllers.DenyFollowRequest(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.WriteResponse()
}

func HandleUnfollowUser(data state.StateHandler) {
	var err error
	err = parsers.ParseUnfollowUserRequest(data)
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	err = controllers.UnfollowUser(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.WriteResponse()
}
