package main

import (
	"os"

	"github.com/tjblackheart/wpdf/service"
)

func main() {
	s := service.New(os.Getenv("APP_HOST"))
	s.Serve()
}
