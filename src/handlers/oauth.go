package handlers

import (
	"errors"
	"forum/src/controllers"
	"forum/src/ferror"
	"forum/src/state"
	"forum/src/utils"
)

func HandleOAuthLoginGoogle(data state.StateHandler) {
	controllers.HandleOAuthLogin(data.(state.StateController), "google")
}

func HandleOAuthLoginGithub(data state.StateHandler) {
	controllers.HandleOAuthLogin(data.(state.StateController), "github")
}

func HandleGoogleCallback(data state.StateHandler) {
	cookieState, err := data.GetRequest().Cookie("__Host-FRMState")
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		data.SetErrorConsume(ferror.ErrorCookieNotFound)
		data.(state.StateController).WriteResponse()
		return
	}
	urlState := data.GetRequest().URL.Query().Get("state")
	if cookieState.Value != urlState {
		err = errors.Join(utils.GetFunctionName(), err)
		data.SetErrorConsume(ferror.ErrorInvalidOAuthState)
		data.(state.StateController).WriteResponse()
		return
	}
	err = controllers.OAuthGoogleCallback(data.(state.StateController))
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	data.(state.StateController).WriteResponse()
}

func HandleGitHubCallback(data state.StateHandler) {
	cookieState, err := data.GetRequest().Cookie("__Host-FRMState")
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		data.SetErrorConsume(ferror.ErrorCookieNotFound)
		data.(state.StateController).WriteResponse()
		return
	}
	urlState := data.GetRequest().URL.Query().Get("state")
	if cookieState.Value != urlState {
		err = errors.Join(utils.GetFunctionName(), ferror.ErrorInvalidOAuthState)
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	err = controllers.OAuthGitHubCallback(data.(state.StateController))
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	data.(state.StateController).WriteResponse()
}
