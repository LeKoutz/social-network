package handlers

import (
	"forum/src/controllers"
	"forum/src/ferror"
	"forum/src/models"
	"forum/src/parsers"
	"forum/src/state"
	"net/http"
)

func HandleShowProfileView(data state.StateHandler) {
	var err error
	data.EditProfile().UserId, err = parsers.ParseProfileId(data)
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	err = controllers.ShowProfile(data.(state.StateController))
	if err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	data.WriteResponse()
}

func HandleEditProfile(data state.StateHandler) {
	req := data.GetRequest()

	if req.Method != http.MethodGet && req.Method != http.MethodPost {
		data.SetErrorConsume(ferror.ErrorMethodNotAllowed).WriteResponse()
		return
	}

	controllerData := data.(state.StateController)

	if err := controllers.ShowEditProfile(controllerData); err != nil {
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}

	if req.Method == http.MethodPost {
		req.Body = http.MaxBytesReader(
			*data.EditResponse(),
			req.Body,
			models.MaxImageSize+(1<<20),
		)

		defer func() {
			profile := data.EditProfile()
			if profile.AvatarFile != nil {
				_ = profile.AvatarFile.Close()
				profile.AvatarFile = nil
			}
			if req.MultipartForm != nil {
				_ = req.MultipartForm.RemoveAll()
			}
		}()

		if err := parsers.ParseProfileEditRequest(data); err != nil {
			data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
			return
		}

		if err := controllers.UpdateProfile(controllerData); err != nil {
			data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
			return
		}
	}
	data.WriteResponse()
}
