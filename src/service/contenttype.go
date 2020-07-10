package service

import (
	"net/http"
	"os"
)

func (s Service) contentType(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	buffer := make([]byte, 512)
	_, err = f.Read(buffer)
	if err != nil {
		return "", err
	}

	t := http.DetectContentType(buffer)
	if err != nil {
		return "", err
	}

	return t, nil
}
