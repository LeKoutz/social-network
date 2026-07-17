package controllers

import (
	"forum/src/models"
)

func getUsers(data models.ResponseStruct) {
	users, err := models.GetAllUsers()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	onlineUsers := Hub.GetOnlineUsers()
	for i := range users {
		users[i].LoggedIn = onlineUsers[users[i].Id]
	}
	data.Users = users
	data.WriteResponse() // no view necessary
}