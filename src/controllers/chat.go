package controllers

import (
	"forum/src/models"
	"forum/src/utils"
	"strings"
)

func showChatHistory(data models.ResponseStruct) {
	id, ok := strings.CutPrefix(data.Request.RequestURI, "/api/chat/")
	if !ok || len(id) == 0 {
		(&models.Error{}).Consume(models.ErrorBadRequest).LogAndRespondError(data.Response, data.User)
		return
	}
	chatUserId, err := utils.StringToInt64(id)
	if err != nil {
		(&models.Error{}).Consume(models.ErrorInvalidChatId).LogAndRespondError(data.Response, data.User)
		return
	}
	messages, err := models.GetChatHistory(data.User.Id, chatUserId)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	data.User.ChatMessages = messages
	data.WriteResponse()
}