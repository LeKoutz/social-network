package handlers

import (
	"forum/src/controllers"
	"forum/src/state"
)

func HandleWs(data state.StateHandler) {
	controllers.ServeWs(data.(state.StateController))
}
