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
	if strings.HasSuffix(data.GetRequest().URL.Path, "/") {
		return "", ferror.ErrorNotFound
	}
	imgURL := filepath.Base(data.GetRequest().URL.Path)
	_, err := os.Stat("./uploads/images/" + imgURL)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", ferror.ErrorNotFound
		} else {
			return "", err
		}
	}
	return imgURL, nil
}
