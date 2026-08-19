package middleware

import (
	"forum/src/state"
)

func AuthMiddleware(data state.State) string {
	for _, cookie := range data.Request.Cookies() {
		if cookie.Name == "__Host-FRMSessionID" {
			return cookie.Value
		}
	}
	return ""
}