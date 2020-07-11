package service

import (
	"errors"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

var validTypes = []string{
	"text/plain",
	"text/html",
}

func (s Service) isSupported(path string) error {
	mime, err := mimetype.DetectFile(path)
	if err != nil {
		return err
	}

	mimeType := mime.String()
	for _, t := range validTypes {
		if strings.Contains(mimeType, t) {
			return nil
		}
	}

	return errors.New("unsupported content-type " + mimeType)
}
