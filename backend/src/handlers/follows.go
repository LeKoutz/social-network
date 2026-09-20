package handlers

import (
	"errors"
	"forum/src/controllers"
	"forum/src/parsers"
	"forum/src/state"
	"forum/src/utils"
)


func HandleUserFollow(data state.StateHandler) {
	var err error
	err = parsers.ParseUserFollowRequest(data)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		data.SetErrorConsume(err).WriteResponse()
		return
	}
	err = controllers.CreateFollowInvitation(data.(state.StateController))
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		data.SetErrorConsume(err).WriteResponse()
		return
	}
	data.WriteResponse()
}
