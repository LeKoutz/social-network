package handlers

import (
	"errors"
	"forum/src/ferror"
	"forum/src/state"
	"forum/src/utils"
	"net/http"
	"os"
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
			err = ferror.ErrorNotFound
		}
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	if stat.IsDir() {
		err = ferror.ErrorCommentPermissionDenied
		data.SetErrorConsume(ferror.ReturnErr(err)).WriteResponse()
		return
	}
	http.ServeFile(*data.EditResponse(), data.GetRequest(), "./public/"+fileURL)
}
