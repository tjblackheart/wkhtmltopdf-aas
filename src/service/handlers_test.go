package service

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleRequest(t *testing.T) {
	s := New(":3000", "/tmp", "/usr/bin/wkhtmltopdf")

	tests := []struct {
		payload  string
		expected string
		status   int
	}{
		{
			`{ "options": "--title Test", "type": "file", "file": "/tmp/test.html" }`,
			`{ "status": "success", "file": "/tmp/test.pdf" }`,
			http.StatusOK,
		},
		{
			`{ "options": "--title Test", "type": "string", "string": "hello world!" }`,
			`{ "status": "success", "file": "/tmp/output.pdf" }`,
			http.StatusOK,
		},
		{
			`{ "options": "--title Test", "type": "url", "url": "https://www.netzwerkorange.de/en" }`,
			`{ "status": "success", "file": "/tmp/www.netzwerkorange.de.pdf" }`,
			http.StatusOK,
		},
		{
			`{ "options": "--title Test", "type": "file", "file": "not_here" }`,
			`{ "error": "Input file not found. The file must reside in the shared dir."}`,
			http.StatusInternalServerError,
		},
		{
			`{ "options": "--title Test", "type": "file" }`,
			`{ "error": "No filename provided." }`,
			http.StatusBadRequest,
		},
		{
			`{ "type": "string", "string": "" }`,
			`{ "error": "No string provided." }`,
			http.StatusBadRequest,
		},
		{
			`{ "type": "file", "string": "something" }`,
			`{ "error": "No filename provided." }`,
			http.StatusBadRequest,
		},
		{
			`{ "type": "unavailable" }`,
			`{ "error": "Type not available." }`,
			http.StatusBadRequest,
		},
		{
			`{}`,
			`{ "error": "Type not available." }`,
			http.StatusBadRequest,
		},
		{
			`{ "type": "url", "url": "inval.id" }`,
			`{ "error": "Error fetching URL." }`,
			http.StatusInternalServerError,
		},
		{
			`{ "type": "url" }`,
			`{ "error": "No URL provided." }`,
			http.StatusBadRequest,
		},
		{
			`invalid_payload`,
			`{ "error": "Invalid payload." }`,
			http.StatusBadRequest,
		},
		{
			``,
			`{ "error": "Invalid payload." }`,
			http.StatusBadRequest,
		},
	}

	ioutil.WriteFile("/tmp/test.html", []byte("Hello world!"), 0644)

	for _, tt := range tests {
		rr := httptest.NewRecorder()
		r := bytes.NewReader([]byte(tt.payload))

		req, err := http.NewRequest("POST", "/", r)
		assert.NoError(t, err)

		http.HandlerFunc(s.handleRequest).ServeHTTP(rr, req)

		assert.Equal(t, tt.status, rr.Code)
		assert.JSONEq(t, tt.expected, rr.Body.String())
	}
}
