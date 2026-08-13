package handlers

import (
	"errors"
	"forum/src/controllers"
	"forum/src/state"
	"forum/src/utils"
	"net/http"
)

func HandleMarkAllNotificationsAsRead(data state.StateHandler) {
	err := controllers.MarkAllNotificationsAsRead(data.(state.StateController))
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
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
