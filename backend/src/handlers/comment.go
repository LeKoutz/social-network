package handlers

import (
	"forum/src/controllers"
	"forum/src/models"
	"forum/src/parsers"
	"forum/src/state"
	"forum/src/ferror"
)

func HandleCommentCreate(data state.StateHandler) {
	var err error
	err = parsers.ParseCreateCommentRequest(data)
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	err = controllers.CommentCreate(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.WriteResponse()
}

func HandleCommentReaction(data state.StateHandler) {
	var err error
	data.EditComment().Id, err = parsers.ParseCommentId(data)
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	err = controllers.CommentReaction(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.WriteResponse()
}

func HandleCommentDelete(data state.StateHandler) {
	var err error
	data.EditComment().Id, err = parsers.ParseCommentId(data)
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	err = controllers.CommentDelete(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.WriteResponse()
}

func HandleCommentEdit(data state.StateHandler) {
	var err error
	err = data.GetRequest().ParseMultipartForm(models.MaxImageSize)
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.EditComment().Id, err = parsers.ParseCommentId(data)
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	err = controllers.CommentEdit(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.WriteResponse()
}
