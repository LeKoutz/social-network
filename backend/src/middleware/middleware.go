package middleware

import (
	"forum/src/controllers"
	"forum/src/state"
)

func AuthMiddleware(data *state.State) {
	var value string
	for _, cookie := range data.Request.Cookies() {
		if cookie.Name == "__Host-FRMSessionID" {
			value = cookie.Value
			break
		}
	}
	if value != "" {
		data.EditUser().SessionId = value
		controllers.GetReturningUser(data)
	}
}
