package handlers

import (
	"forum/src/state"
)

func HandleShowCategories(data state.StateHandler) {
	data.EditCategories()
}
