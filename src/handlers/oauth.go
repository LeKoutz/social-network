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
		data.SetErrorConsume(ferror.ErrorCookieNotFound)
		data.(state.StateController).WriteResponse()
		return
	}
	urlState := data.GetRequest().URL.Query().Get("state")
	if cookieState.Value != urlState {
		data.SetErrorConsume(ferror.ErrorInvalidOAuthState)
		data.(state.StateController).WriteResponse()
		return
	}
	err = controllers.OAuthGoogleCallback(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	data.(state.StateController).WriteResponse()
}

func HandleGitHubCallback(data state.StateHandler) {
	cookieState, err := data.GetRequest().Cookie("__Host-FRMState")
	if err != nil {
		data.SetErrorConsume(ferror.ErrorCookieNotFound)
		data.(state.StateController).WriteResponse()
		return
	}
	urlState := data.GetRequest().URL.Query().Get("state")
	if cookieState.Value != urlState {
		data.SetErrorConsume(ferror.ErrorInvalidOAuthState)
		data.(state.StateController).WriteResponse()
		return
	}
	err = controllers.OAuthGitHubCallback(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	data.(state.StateController).WriteResponse()
}
