package service

import (
	"errors"
	"net/http"
	"os"
	"strings"
)

func (s Service) contentType(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	buffer := make([]byte, 512)
	_, err = f.Read(buffer)
	if err != nil {
		return err
	}

	t := http.DetectContentType(buffer)
	if err != nil {
		return err
	}

	if !strings.Contains(t, "text/html") || !strings.Contains(t, "text/plain") {
		return errors.New("unsupported content-type: " + t)
	}

	return nil
}
