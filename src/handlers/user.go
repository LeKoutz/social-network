package handlers

import (
	"errors"
	"forum/src/controllers"
	"forum/src/ferror"
	"forum/src/state"
	"forum/src/utils"
	"net/http"
)

func HandleUserLogin(data state.StateHandler) {
	if data.GetUser().LoggedIn {
		data.SetErrorConsume(ferror.ErrorAlreadyLoggedIn)
		data.(state.StateController).WriteResponse()
		return
	}
	switch data.GetRequest().Method {
	case http.MethodPost:
		err := controllers.AttemptLogin(data.(state.StateController))
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			data.SetErrorConsume(err)
			data.(state.StateController).WriteResponse()
			return
		}
		// http.Redirect(*data.EditResponse(), data.GetRequest(), "/", http.StatusSeeOther)
		data.(state.StateController).WriteResponse()
		return
	case http.MethodGet:
		data.(state.StateController).WriteResponse()
		return
	default:
		data.SetErrorConsume(ferror.ErrorMethodNotAllowed)
		data.(state.StateController).WriteResponse()
		return
	}
}

func HandleUserLogout(data state.StateHandler) {
	var err error
	err = controllers.UserLogout(data.(state.StateController))
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	data.(state.StateController).WriteResponse()
}

func HandleShowUserPosts(data state.StateHandler) {
	var err error
	err = controllers.GetUserPosts(data.(state.StateController))
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	data.(state.StateController).WriteResponse()
}

func HandleShowUserLikedPosts(data state.StateHandler) {
	var err error
	err = controllers.GetUserLikedPosts(data.(state.StateController))
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	data.(state.StateController).WriteResponse()
}

func HandleShowUserView(data state.StateHandler) {
	data.(state.StateController).WriteResponse()
}

func HandleShowUserActivity(data state.StateHandler) {
	err := controllers.GetUserActivity(data.(state.StateController))
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		data.SetErrorConsume(ferror.ErrorInternalServerError)
		data.(state.StateController).WriteResponse()
		return
	}
	data.(state.StateController).WriteResponse()
}

func HandleUserRegister(data state.StateHandler) {
	if data.GetUser().LoggedIn {
		data.SetErrorConsume(ferror.ErrorAlreadyLoggedIn)
		data.(state.StateController).WriteResponse()
		return
	}
	switch data.GetRequest().Method {
	case http.MethodPost:
		if len(data.GetRequest().FormValue("username")) == 0 ||
			len(data.GetRequest().FormValue("email")) == 0 ||
			len(data.GetRequest().FormValue("first_name")) == 0 ||
			len(data.GetRequest().FormValue("last_name")) == 0 ||
			len(data.GetRequest().FormValue("age")) == 0 ||
			len(data.GetRequest().FormValue("gender")) == 0 ||
			len(data.GetRequest().FormValue("password1")) == 0 ||
			len(data.GetRequest().FormValue("password2")) == 0 {
			data.SetErrorConsume(ferror.ErrorBadRequest)
			data.(state.StateController).WriteResponse()
			return
		}
		data.EditUser().Username = data.GetRequest().FormValue("username")
		data.EditUser().Email = data.GetRequest().FormValue("email")
		controllers.AttemptRegister(data.(state.StateController))
		return
	default:
		data.SetErrorConsume(ferror.ErrorMethodNotAllowed)
		data.(state.StateController).WriteResponse()
	}
}

func HandleGetUsers(data state.StateHandler) {
	var err error
	err = controllers.GetUsersForPanel(data.(state.StateController))
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	controllers.HubOnlineUsers(data.(state.StateController))
	data.(state.StateController).WriteResponse()
}
