package handlers

import (
	"forum/src/controllers"
	"forum/src/ferror"
	"forum/src/state"
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
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	urlState := data.GetRequest().URL.Query().Get("state")
	if cookieState.Value != urlState {
		err = ferror.ErrorInvalidOAuthState
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	err = controllers.OAuthGoogleCallback(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.WriteResponse()
}

func HandleGitHubCallback(data state.StateHandler) {
	cookieState, err := data.GetRequest().Cookie("__Host-FRMState")
	if err != nil {
		err = ferror.ErrorCookieNotFound
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	urlState := data.GetRequest().URL.Query().Get("state")
	if cookieState.Value != urlState {
		err = ferror.ErrorInvalidOAuthState
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	err = controllers.OAuthGitHubCallback(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.WriteResponse()
}
