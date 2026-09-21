package handlers

import (
	"forum/src/controllers"
	"forum/src/state"
	"forum/src/ferror"
)

func HandleMarkAllNotificationsAsRead(data state.StateHandler) {
	err := controllers.MarkAllNotificationsAsRead(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.WriteResponse()
}
