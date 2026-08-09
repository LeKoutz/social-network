package handlers

import (
	"forum/src/controllers"
	"forum/src/state"
	"net/http"
)

func HandleMarkAllNotificationsAsRead(data state.StateHandler) {
	err := controllers.MarkAllNotificationsAsRead(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	previousURL := data.GetRequest().Referer()
	if previousURL == "" {
		previousURL = "/"
	}
	http.Redirect(*data.EditResponse(), data.GetRequest(), previousURL, http.StatusSeeOther)
}
