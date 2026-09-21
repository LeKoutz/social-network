package handlers

import (
	"forum/src/controllers"
	"forum/src/parsers"
	"forum/src/state"
	"forum/src/ferror"
)

func HandleServeUnreadMessages(data state.StateHandler) {
	err := controllers.ServeUnreadMessages(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.WriteResponse()
}

func HandleShowChatHistory(data state.StateHandler) {
	recipientId, offset, err := parsers.ParseChatId(data)
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.EditChatMessage(0).RecipientId = recipientId
	data.SetChatOffset(offset)
	err = controllers.ShowChatHistory(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.WriteResponse()
}
