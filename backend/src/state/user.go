package state

import "forum/src/models"

func (r *State) InitUser() *State {
	r.User = models.GetGuestUser()
	return r
}

func (r *State) EditUser() *models.UserType {
	return &r.User
}

func (r State) GetUser() models.UserType {
	return r.User
}

func (r *State) SetUser(user models.UserType) *State {
	r.User = user
	return r
}

func (r *State) EditUsers() *models.UsersType {
	return &r.Users
}

func (r *State) GetUsers() models.UsersType {
	return r.Users
}

func (r *State) SetUsers(users models.UsersType) *State {
	r.Users = users
	return r
}
