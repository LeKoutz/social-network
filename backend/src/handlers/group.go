package handlers

import (
	"errors"
	"forum/src/controllers"
	"forum/src/ferror"
	"forum/src/parsers"
	"forum/src/state"
	"forum/src/utils"
	"net/http"
)

func HandleCreateGroup(data state.StateHandler) {
	if data.GetRequest().Method != http.MethodPost {
		data.SetErrorConsume(ferror.ErrorMethodNotAllowed)
		data.WriteResponse()
		return
	}
	var err error
	err = parsers.ParseCreateGroupRequest(data)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		data.SetErrorConsume(err).WriteResponse()
		return
	}
	err = controllers.CreateGroup(data.(state.StateController))
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		data.SetErrorConsume(err).WriteResponse()
		return
	}
	data.WriteResponse()
}