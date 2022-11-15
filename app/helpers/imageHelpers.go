package helpers

import (
	"errors"
	"strings"
)

var InvalidImgExtension = errors.New("File-Name should have an extension .png .jpeg or .gif ")

func ValidateImgFileNameHeader(fileName string) error {
	if fileName == "" {
		return errors.New("File-Name header should be specified")
	}

	if !strings.Contains(fileName, ".") {
		return nil
	}

	extS := strings.Split(fileName, ".")
	ext := strings.ToLower(extS[len(extS)-1])
	if ext != "png" && ext != "jpeg" && ext != "gif" {
		return InvalidImgExtension
	}

	return nil
}
