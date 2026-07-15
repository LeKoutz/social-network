package controllers

import (
	"forum/src/models"
	"forum/src/utils"
	"strings"
)

func showChatHistory(data models.ResponseStruct) {
	uri, ok := strings.CutPrefix(data.Request.RequestURI, "/api/chat/")
	if !ok {
		(&models.Error{}).Consume(models.ErrorBadRequest).LogAndRespondError(data.Response, data.User)
		return
	}
	id, _, _ := strings.Cut(uri, "?")
	if len(id) == 0 {
		(&models.Error{}).Consume(models.ErrorBadRequest).LogAndRespondError(data.Response, data.User)
		return
	}
	chatUserId, err := utils.StringToInt64(id)
	if err != nil {
		(&models.Error{}).Consume(models.ErrorInvalidChatId).LogAndRespondError(data.Response, data.User)
		return
	}
	offset, _ := utils.StringToInt64(data.Request.URL.Query().Get("offset"))
	messages, err := models.GetChatHistory(data.User.Id, chatUserId, offset)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	for i := range messages {
		if messages[i].RecipientId == data.User.Id {
			messages[i].MarkAsRead()
		}
	}
	data.User.ChatMessages = messages
	data.WriteResponse()
}

func serveUnreadMessages(data models.ResponseStruct) {
	messages, err := models.GetUnreadMessageIds(data.User.Id)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	data.User.ChatMessages = messages
	data.WriteResponse()
}