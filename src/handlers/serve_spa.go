package handlers

import (
	"errors"
	"forum/src/ferror"
	"forum/src/state"
	"forum/src/utils"
	"net/http"
	"os"
	// "path/filepath"
	// "strings"
)

func HandleServeSPA(data state.StateHandler) {
	utils.LogInfo(data.GetRequest().URL.Path)
	fileURL := data.GetRequest().URL.Path
	if fileURL == "" || fileURL == "/" {
		fileURL = "index.html"
	}
	stat, err := os.Stat("./public/" + fileURL)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			data.SetErrorConsume(ferror.ErrorNotFound)
			data.(state.StateController).WriteResponse()
			return
		}
		data.SetErrorConsume(err)
		data.(state.StateController).WriteResponse()
		return
	}
	if stat.IsDir() {
		data.SetErrorConsume(ferror.ErrorPermissionDenied)
		data.(state.StateController).WriteResponse()
		return
	}
	http.ServeFile(*data.EditResponse(), data.GetRequest(), "./public/"+fileURL)
}
