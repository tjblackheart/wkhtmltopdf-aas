package service

import (
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	s := New(":3000", "/tmp", "/usr/bin/wkhtmltopdf")
	assert.IsType(t, &Service{}, s)
	assert.Equal(t, &Service{
		hostname: ":3000",
		shared:   "/tmp",
		binary:   "/usr/bin/wkhtmltopdf",
	}, s)
}

func TestRouter(t *testing.T) {
	s := New(":3000", "/tmp", "/usr/bin/wkhtmltopdf")
	assert.IsType(t, s.router(), &mux.Router{})
}
