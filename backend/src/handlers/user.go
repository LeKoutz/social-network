package handlers

import (
	"forum/src/controllers"
	"forum/src/ferror"
	"forum/src/parsers"
	"forum/src/state"
)

func HandleUserLogin(data state.StateHandler) {
	var err error
	if err = parsers.ParseLoginForm(data); err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	err = controllers.AttemptLogin(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.WriteResponse()
	return
}

func HandleUserLogout(data state.StateHandler) {
	var err error
	err = controllers.UserLogout(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.WriteResponse()
}

func HandleShowUserPosts(data state.StateHandler) {
	var err error
	err = controllers.GetUserPosts(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.WriteResponse()
}

func HandleShowUserLikedPosts(data state.StateHandler) {
	var err error
	err = controllers.GetUserLikedPosts(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.WriteResponse()
}

func HandleUserRegister(data state.StateHandler) {
	err := parsers.ParseRegistrationForm(data)
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	err = controllers.AttemptRegister(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.WriteResponse()
}

func HandleGetUsers(data state.StateHandler) {
	var err error
	err = controllers.GetUsersForPanel(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	controllers.HubOnlineUsers(data.(state.StateController))
	data.WriteResponse()
}
