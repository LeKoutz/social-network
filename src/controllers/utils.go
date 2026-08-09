package controllers

import (
	"fmt"
	"forum/src/state"
)

func GetRedirectLinkToCommentOfPost(data state.StateController) string {
	return fmt.Sprintf("/post/view/%d#comment-%d", data.GetPost().Id, data.GetComment().Id)
}

func GetRedirectLinkToPost(data state.StateController) string {
	return fmt.Sprintf("/post/view/%d", data.GetPost().Id)
}
