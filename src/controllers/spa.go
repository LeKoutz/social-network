package controllers

import (
	"errors"
	"forum/src/models"
	"forum/src/utils"
	"net/http"
	"os"
	// "path/filepath"
	// "strings"
)

func serveSPA(data models.ResponseStruct) {
	utils.LogInfo(data.Request.URL.Path)
	// if strings.HasPrefix(data.Request.URL.Path, "/") {
	// 	(&models.Error{}).Consume(models.ErrorNotFound).LogAndRespondError(data.Response, data.User)
	// 	return
	// }
	fileURL := data.Request.URL.Path
	_, err := os.Stat("./public/" + fileURL)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			(&models.Error{}).Consume(models.ErrorNotFound).LogAndRespondError(data.Response, data.User)
			return
		} else {
			(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
			return
		}
	}
	http.ServeFile(data.Response, data.Request, "./public/"+fileURL)
}

