package controllers

import (
	"forum/src/models"
)
// GetUsersForPanel retrieves all users for the panel, excluding the current user
func getUsers(data models.ResponseStruct) {
	// Get all users for the panel
	users, err := models.GetUsersForPanel(data.User.Id)
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