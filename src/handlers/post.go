package handlers

import (
	"forum/src/controllers"
	"forum/src/ferror"
	"forum/src/models"
	"forum/src/state"
	"forum/src/utils"
	"net/http"
)

func HandleShowPost(data state.StateHandler) {
	var err error
	err = controllers.GetPost(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	data.(state.StateController).WriteResponse()
}

func HandlePostCreateGet(data state.StateHandler) {
	err := (data.(state.StateController)).EditCategories().GetAll()
	if err != nil {
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	data.(state.StateController).WriteResponse()
}

func HandlePostCreatePost(data state.StateHandler) {
	var err error
	err = data.GetRequest().ParseMultipartForm(models.MaxImageSize)
	if err != nil {
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	err = controllers.CreatePost(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	data.(state.StateController).WriteResponse()
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
		data.(state.StateController).WriteResponse()
		return
	}
}

func HandlePostEdit(data state.StateHandler) {
	var err error
	switch data.GetRequest().Method {
	case http.MethodGet:
		err = controllers.ShowEditPost(data.(state.StateController))
		if err != nil {
			data.SetErrorConsume(err)
			data.(state.StateController).WriteResponse()
			return
		}
		data.(state.StateController).WriteResponse()
		return
	case http.MethodPost:
		err = controllers.UpdatePost(data.(state.StateController))
		if err != nil {
			data.SetErrorConsume(err)
			data.(state.StateController).WriteResponse()
			return
		}
		utils.LogDebug(data)
		http.Redirect(*data.EditResponse(), data.GetRequest(), data.GetRedirect(), http.StatusSeeOther)
		return
	default:
		data.SetErrorConsume(ferror.ErrorMethodNotAllowed)
		data.(state.StateController).WriteResponse()
	}
}

func HandlePostReaction(data state.StateHandler) {
	var err error
	err = controllers.PostReaction(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	HandleShowPost(data)
}

func HandlePostDelete(data state.StateHandler) {
	var err error
	var post models.PostType
	post.Id, err = controllers.ParsePostId(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ErrorInvalidPostId)
		data.(state.StateController).WriteResponse()
		return
	}
	err = controllers.RemovePost(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	http.Redirect(*data.EditResponse(), data.GetRequest(), data.GetRedirect(), http.StatusSeeOther)
}
