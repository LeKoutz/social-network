package controllers

import (
	"forum/src/models"
	"net/http"
)

func markAllNotificationsAsRead(data models.ResponseStruct) {
	if !data.User.LoggedIn {
		data.SetErrorConsume(models.ErrorUnauthorizedAction).WriteResponse()
		return
	}
	err := data.User.MarkAllNotificationsAsRead()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	for i := range data.User.Notifications {
		data.User.Notifications[i].Read = true
	}
	previousURL := data.Request.Referer()
	if previousURL == "" {
		previousURL = "/"
	}
	http.Redirect(data.Response, data.Request, previousURL, http.StatusSeeOther)
}
