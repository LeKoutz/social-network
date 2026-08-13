package handlers

import (
	"forum/src/controllers"
	"errors"
	"forum/src/parsers"
	"forum/src/state"
	"forum/src/utils"
	"net/http"
)

func HandleCommentCreate(data state.StateHandler) {
	var err error
	err = parsers.ParseCreateCommentRequest(data)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	err = controllers.CommentCreate(data.(state.StateController))
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	http.Redirect(*data.EditResponse(), data.GetRequest(), data.GetRedirect(), http.StatusSeeOther)
}

func HandleCommentReaction(data state.StateHandler) {
	var err error
	data.EditComment().Id, err = parsers.ParseCommentId(data)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	err = controllers.CommentReaction(data.(state.StateController))
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	http.Redirect(*data.EditResponse(), data.GetRequest(), data.GetRedirect(), http.StatusSeeOther)
}

func HandleCommentDelete(data state.StateHandler) {
	var err error
	data.EditComment().Id, err = parsers.ParseCommentId(data)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	err = controllers.CommentDelete(data.(state.StateController))
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	http.Redirect(*data.EditResponse(), data.GetRequest(), data.GetRedirect(), http.StatusSeeOther)
}

func HandleCommentEdit(data state.StateHandler) {
	var err error
	data.EditComment().Id, err = parsers.ParseCommentId(data)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	err = controllers.CommentEdit(data.(state.StateController))
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	data.(state.StateController).WriteResponse()
}
