package controllers

import (
	"errors"
	"forum/src/ferror"
	"forum/src/state"
	"os"
	"path/filepath"
	"strings"
)

func HandleImages(data state.StateController) (string, error) {
	var err error
	if strings.HasSuffix(data.GetRequest().URL.Path, "/") {
		err = ferror.ErrorNotFound
		return "", ferror.ReturnErr(err)
	}
	imgURL := filepath.Base(data.GetRequest().URL.Path)
	_, err = os.Stat("./uploads/images/" + imgURL)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = ferror.ErrorNotFound
		}
		return "", ferror.ReturnErr(err)
	}
	return imgURL, nil
}
