package controllers

import (
	"errors"
	"forum/src/state"
	"forum/src/utils"
)

func MarkAllNotificationsAsRead(data state.StateController) error {
	err := data.EditUser().MarkAllNotificationsAsRead()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	for i := range data.GetUser().Notifications {
		data.EditUser().Notifications[i].Read = true
	}
	data.EditUser().UnreadNotificationsCount = 0
	return nil
}
