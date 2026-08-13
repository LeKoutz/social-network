package handlers

import (
	"forum/src/controllers"
	"forum/src/ferror"
	"forum/src/state"
	"forum/src/utils"
	"strings"
)

func parseChatId(data state.StateHandler) (id, offset int64, err error) {
	uri, ok := strings.CutPrefix(data.GetRequest().RequestURI, "/api/chat/")
	if !ok {
		return 0, 0, ferror.ErrorBadRequest
	}
	idStr, offsetStr, found := strings.Cut(uri, "?offset=")
	if !found {
		return 0, 0, ferror.ErrorBadRequest
	}
	id, err = utils.StringToInt64(idStr)
	if err != nil {
		return 0, 0, ferror.ErrorInvalidChatId
	}
	offset, err = utils.StringToInt64(offsetStr)
	if err != nil {
		return 0, 0, ferror.ErrorBadRequest
	}
	return id, offset, err
}

func HandleServeUnreadMessages(data state.StateHandler) {
	err := controllers.ServeUnreadMessages(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	data.(state.StateController).WriteResponse()
}

func HandleShowChatHistory(data state.StateHandler) {
	recipientId, offset, err := parseChatId(data)
	if err != nil {
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	data.EditChatMessage(0).RecipientId = recipientId
	data.SetChatOffset(offset)
	err = controllers.ShowChatHistory(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	data.(state.StateController).WriteResponse()
}
