package handlers

import (
	"forum/src/controllers"
	"forum/src/state"
	"net/http"
)

func HandleCommentCreate(data state.StateHandler) {
	var err error
	err = controllers.CommentCreate(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	http.Redirect(*data.EditResponse(), data.GetRequest(), data.GetRedirect(), http.StatusSeeOther)
}

func HandleCommentReaction(data state.StateHandler) {
	var err error
	err = controllers.CommentReaction(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	http.Redirect(*data.EditResponse(), data.GetRequest(), data.GetRedirect(), http.StatusSeeOther)
}

func HandleCommentDelete(data state.StateHandler) {
	var err error
	err = controllers.CommentDelete(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	http.Redirect(*data.EditResponse(), data.GetRequest(), data.GetRedirect(), http.StatusSeeOther)
}

func HandleCommentEdit(data state.StateHandler) {
	var err error
	err = controllers.CommentEdit(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	data.(state.StateController).WriteResponse()
}
