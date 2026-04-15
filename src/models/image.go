package models

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"forum/src/utils"

	"github.com/gofrs/uuid"
)

const MaxImageSize = 20 << 20

func ValidateImage(reader *bytes.Reader) error {
	buf := make([]byte, 512)
	n, err := reader.Read(buf)
	if err != nil && err != io.EOF {
		return err
	}
	buf = buf[:n]

	if reader.Len() > MaxImageSize {
		return ErrorImageTooBig
	}

	if !isValidImageType(buf) {
		return ErrorInvalidImageType
	}

	return nil
}

func isValidImageType(buf []byte) bool {
	if len(buf) < 4 {
		return false
	}

	if bytes.HasPrefix(buf, []byte{0xFF, 0xD8, 0xFF}) {
		return true
	}
	if bytes.HasPrefix(buf, []byte{0x89, 0x50, 0x4E, 0x47}) {
		return true
	}
	if bytes.HasPrefix(buf, []byte{0x47, 0x49, 0x46, 0x38}) {
		return true
	}

	return false
}

func SaveImage(file multipart.File, postId int64, uploadDir string) (string, error) {
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return "", err
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return "", err
	}

	if len(fileBytes) > MaxImageSize {
		return "", ErrorImageTooBig
	}

	if !isValidImageType(fileBytes[:512]) {
		return "", ErrorInvalidImageType
	}

	ext := getImageExtension(fileBytes)
	filename := fmt.Sprintf("%s%s", uuid.Must(uuid.NewV4()).String(), ext)

	dir := filepath.Join(uploadDir, "static", "images")
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return "", err
	}

	dst := filepath.Join(dir, filename)
	out, err := os.Create(dst)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return "", err
	}
	defer out.Close()

	_, err = out.Write(fileBytes)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return "", err
	}

	return filepath.Join("static", "images", filename), nil
}

func getImageExtension(buf []byte) string {
	if bytes.HasPrefix(buf, []byte{0xFF, 0xD8, 0xFF}) {
		return ".jpg"
	}
	if bytes.HasPrefix(buf, []byte{0x89, 0x50, 0x4E, 0x47}) {
		return ".png"
	}
	if bytes.HasPrefix(buf, []byte{0x47, 0x49, 0x46, 0x38}) {
		return ".gif"
	}
	return ""
}
