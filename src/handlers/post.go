package handlers

import (
	"errors"
	"forum/src/controllers"
	"forum/src/ferror"
	"forum/src/models"
	"forum/src/parsers"
	"forum/src/state"
	"forum/src/utils"
	"net/http"
)

func HandleShowPost(data state.StateHandler) {
	var err error
	data.EditPost().Id, err = parsers.ParsePostId(data)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		data.SetErrorConsume(err).WriteResponse()
		return
	}
	err = controllers.GetPost(data.(state.StateController))
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		data.SetErrorConsume(err).WriteResponse()
		return
	}
	data.WriteResponse()
}

func HandlePostCreateGet(data state.StateHandler) {
	err := (data.(state.StateController)).EditCategories().GetAll()
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		data.SetErrorConsume(err).WriteResponse()
		return
	}
	data.WriteResponse()
}

func HandlePostCreatePost(data state.StateHandler) {
	var err error
	err = parsers.ParseCreatePostRequest(data)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		data.SetErrorConsume(err).WriteResponse()
		return
	}
	err = data.GetRequest().ParseMultipartForm(models.MaxImageSize)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		data.SetErrorConsume(err).WriteResponse()
		return
	}
	err = controllers.CreatePost(data.(state.StateController))
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		data.SetErrorConsume(err).WriteResponse()
		return
	}
	data.WriteResponse()
	// http.Redirect(*data.EditResponse(), data.GetRequest(), data.GetRedirect(), http.StatusSeeOther)
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
		err = errors.Join(utils.GetFunctionName(), err)
		data.SetErrorConsume(err).WriteResponse()
		return
	}
	switch data.GetRequest().Method {
	case http.MethodGet:
		err = controllers.ShowEditPost(data.(state.StateController))
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			data.SetErrorConsume(err).WriteResponse()
			return
		}
		data.WriteResponse()
		return
	case http.MethodPost:
		err = parsers.ParseCreatePostRequest(data)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			data.SetErrorConsume(err).WriteResponse()
			return
		}
		err = controllers.UpdatePost(data.(state.StateController))
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			data.SetErrorConsume(err).WriteResponse()
			return
		}
		utils.LogDebug(data)
		http.Redirect(*data.EditResponse(), data.GetRequest(), data.GetRedirect(), http.StatusSeeOther)
		return
	default:
		data.SetErrorConsume(ferror.ErrorMethodNotAllowed).WriteResponse()
	}
}

func HandlePostReaction(data state.StateHandler) {
	var err error
	data.EditPost().Id, err = parsers.ParsePostId(data)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		data.SetErrorConsume(err).WriteResponse()
		return
	}
	err = controllers.PostReaction(data.(state.StateController))
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		data.SetErrorConsume(err).WriteResponse()
		return
	}
	HandleShowPost(data)
}

func HandlePostDelete(data state.StateHandler) {
	var err error
	var post models.PostType
	post.Id, err = parsers.ParsePostId(data)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		data.SetErrorConsume(ferror.ErrorInvalidPostId).WriteResponse()
		return
	}
	data.(state.StateController).SetPost(post)
	err = controllers.RemovePost(data.(state.StateController))
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		data.SetErrorConsume(err).WriteResponse()
		return
	}
	http.Redirect(*data.EditResponse(), data.GetRequest(), data.GetRedirect(), http.StatusSeeOther)
}
