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
		err = ferror.ErrorCookieNotFound
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		data.SetErrorConsume(err).WriteResponse()
		return
	}
	urlState := data.GetRequest().URL.Query().Get("state")
	if cookieState.Value != urlState {
		err = ferror.ErrorInvalidOAuthState
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		data.SetErrorConsume(err).WriteResponse()
		return
	}
	err = controllers.OAuthGoogleCallback(data.(state.StateController))
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		data.SetErrorConsume(err).WriteResponse()
		return
	}
	data.WriteResponse()
}

func HandleGitHubCallback(data state.StateHandler) {
	cookieState, err := data.GetRequest().Cookie("__Host-FRMState")
	if err != nil {
		err = ferror.ErrorCookieNotFound
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		data.SetErrorConsume(err).WriteResponse()
		return
	}
	urlState := data.GetRequest().URL.Query().Get("state")
	if cookieState.Value != urlState {
		err = ferror.ErrorInvalidOAuthState
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		data.SetErrorConsume(err).WriteResponse()
		return
	}
	err = controllers.OAuthGitHubCallback(data.(state.StateController))
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		data.SetErrorConsume(err).WriteResponse()
		return
	}
	data.WriteResponse()
}
