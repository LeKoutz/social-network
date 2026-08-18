package controllers

import (
	"errors"
	"forum/src/models"
	"forum/src/state"
	"forum/src/utils"
)

func GetUsersForPanel(data state.StateController) error {
	err := data.EditUsers().GetUsersForPanel(data.GetUser().Id)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}

func HubOnlineUsers(data state.StateController) {
	onlineUsers := models.MainHub.GetOnlineUsers()
	users := data.GetUsers()
	for i := range users {
		users[i].LoggedIn = onlineUsers[users[i].Id]
	}
	data.SetUsers(users)
}
