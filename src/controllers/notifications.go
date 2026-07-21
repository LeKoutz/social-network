package controllers

import (
	"forum/src/models"
)

func markAllNotificationsAsRead(data models.ResponseStruct) {
	err := data.User.MarkAllNotificationsAsRead()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	for i := range data.User.Notifications {
		data.User.Notifications[i].Read = true
	}
	data.User.UnreadNotificationsCount = 0
	data.WriteResponse()
}
