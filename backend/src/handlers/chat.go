package handlers

import (
	"errors"
	"forum/src/controllers"
	"forum/src/parsers"
	"forum/src/state"
	"forum/src/utils"
)

func HandleServeUnreadMessages(data state.StateHandler) {
	err := controllers.ServeUnreadMessages(data.(state.StateController))
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		data.SetErrorConsume(err).WriteResponse()
		return
	}
	data.WriteResponse()
}

func HandleShowChatHistory(data state.StateHandler) {
	recipientId, offset, err := parsers.ParseChatId(data)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		data.SetErrorConsume(err).WriteResponse()
		return
	}
	data.EditChatMessage(0).RecipientId = recipientId
	data.SetChatOffset(offset)
	err = controllers.ShowChatHistory(data.(state.StateController))
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		data.SetErrorConsume(err).WriteResponse()
		return
	}
	data.WriteResponse()
}
