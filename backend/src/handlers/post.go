package handlers

import (
	"forum/src/controllers"
	"forum/src/ferror"
	"forum/src/models"
	"forum/src/parsers"
	"forum/src/state"
	"net/http"
)

func HandleShowPost(data state.StateHandler) {
	var err error
	data.EditPost().Id, err = parsers.ParsePostId(data)
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	err = controllers.GetPost(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.WriteResponse()
}

func HandlePostCreateGet(data state.StateHandler) {
	err := (data.(state.StateController)).EditCategories().GetAll()
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.WriteResponse()
}

func HandlePostCreatePost(data state.StateHandler) {
	var err error
	err = parsers.ParseCreatePostRequest(data)
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	err = data.GetRequest().ParseMultipartForm(models.MaxImageSize)
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	err = controllers.CreatePost(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.WriteResponse()
}

func HandlePostCreate(data state.StateHandler) {
	switch data.GetRequest().Method {
	case http.MethodGet:
		HandlePostCreateGet(data)
		return
	case http.MethodPost:
		HandlePostCreatePost(data)
		return
	default:
		data.SetErrorConsume(ferror.ErrorMethodNotAllowed)
		data.WriteResponse()
		return
	}
}

func HandlePostEdit(data state.StateHandler) {
	var err error
	data.EditPost().Id, err = parsers.ParsePostId(data)
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	switch data.GetRequest().Method {
	case http.MethodGet:
		err = controllers.ShowEditPost(data.(state.StateController))
		if err != nil {
			data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
			return
		}
		data.WriteResponse()
		return
	case http.MethodPost:
		err = parsers.ParseCreatePostRequest(data)
		if err != nil {
			data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
			return
		}
		err = controllers.UpdatePost(data.(state.StateController))
		if err != nil {
			data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
			return
		}
		data.WriteResponse()
		return
	default:
		err = ferror.ErrorMethodNotAllowed
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
	}
}

func HandlePostReaction(data state.StateHandler) {
	var err error
	data.EditPost().Id, err = parsers.ParsePostId(data)
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	err = controllers.PostReaction(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	HandleShowPost(data)
}

func HandlePostDelete(data state.StateHandler) {
	var err error
	var post models.PostType
	post.Id, err = parsers.ParsePostId(data)
	if err != nil {
		err = ferror.ErrorInvalidPostId
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.(state.StateController).SetPost(post)
	err = controllers.RemovePost(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.WriteResponse()
}
