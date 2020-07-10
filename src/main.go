package main

import (
	"os"

	"github.com/tjblackheart/wpdf/service"
)

func main() {
	host, shared, binary := setup()
	s := service.New(host, shared, binary)
	s.Serve()
}

func setup() (string, string, string) {
	host := os.Getenv("APP_HOST")
	if host == "" {
		host = ":3000"
	}

	shared := os.Getenv("APP_SHARED")
	if shared == "" {
		shared = "/app/shared"
	}

	binary := os.Getenv("APP_BINARY")
	if binary == "" {
		binary = "/usr/local/bin/wkhtmltopdf"
	}

	return host, shared, binary
}
