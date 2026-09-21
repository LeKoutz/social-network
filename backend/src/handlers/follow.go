package handlers

import (
	"forum/src/controllers"
	"forum/src/ferror"
	"forum/src/parsers"
	"forum/src/state"
)

func HandleUserFollow(data state.StateHandler) {
	var err error
	err = parsers.ParseUserFollowRequest(data)
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	err = controllers.CreateFollowInvitation(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.WriteResponse()
}

func HandleUserUnFollow(data state.StateHandler) {
    var err error
    err = parsers.ParseUserFollowRequest(data)
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
