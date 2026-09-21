package controllers

import (
	"forum/src/state"
	"forum/src/ferror"
)

func MarkAllNotificationsAsRead(data state.StateController) error {
	err := data.EditUser().MarkAllNotificationsAsRead()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	for i := range data.GetUser().Notifications {
		data.EditUser().Notifications[i].Read = true
	}
	data.EditUser().UnreadNotificationsCount = 0
	return nil
}
