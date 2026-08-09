package controllers

import (
	"forum/src/db"
	"forum/src/ferror"
	"forum/src/models"
	"forum/src/state"
	"forum/src/utils"
	"strings"
)

func showChatHistory(data state.StateHandler) {
	uri, ok := strings.CutPrefix(data.GetRequest().RequestURI, "/api/chat/")
	if !ok {
		data.SetErrorConsume(ferror.ErrorBadRequest)
		data.(state.StateController).WriteResponse()
		return
	}
	id, _, _ := strings.Cut(uri, "?")
	if len(id) == 0 {
		data.SetErrorConsume(ferror.ErrorBadRequest)
		data.(state.StateController).WriteResponse()
		return
	}
	chatUserId, err := utils.StringToInt64(id)
	if err != nil {
		data.SetErrorConsume(ferror.ErrorInvalidChatId)
		data.(state.StateController).WriteResponse()
		return
	}
	if chatUserId == data.GetUser().Id {
		data.SetErrorConsume(ferror.ErrorNotFound)
		data.(state.StateController).WriteResponse()
		return
	}
	offset, _ := utils.StringToInt64(data.GetRequest().URL.Query().Get("offset"))
	messages, err := db.GetChatHistory(data.GetUser().Id, chatUserId, offset)
	if err != nil {
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	for i := range messages {
		if messages[i].RecipientId == data.GetUser().Id {
			messages[i].MarkAsRead()
		}
	}
	var m models.ChatMessagesType
	m.ConvertFromRow(&messages)
	data.EditUser().ChatMessages = m
	data.(state.StateController).WriteResponse()
}

func serveUnreadMessages(data state.StateHandler) {
	messages, err := db.GetUnreadMessageIds(data.GetUser().Id)
	if err != nil {
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	var m models.ChatMessagesType
	m.ConvertFromRow(&messages)
	data.EditUser().ChatMessages = m
	data.(state.StateController).WriteResponse()
}
