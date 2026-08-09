package controllers

import (
	"errors"
	"forum/src/state"
	"forum/src/utils"
)

func MarkAllNotificationsAsRead(data state.StateController) error {
	err := data.EditUser().MarkAllNotificationsAsRead()
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	for i := range data.GetUser().Notifications {
		data.EditUser().Notifications[i].Read = true
	}
	data.EditUser().UnreadNotificationsCount = 0
	data.WriteResponse()
	return nil
}
