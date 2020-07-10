package service

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleRequest(t *testing.T) {
	s := New(":3000", "/tmp", "/usr/bin/wkhtmltopdf")

	tests := []struct {
		payload  string
		expected string
		status   int
	}{
		{
			`{"options":"--title Test","type":"file","file":"/tmp/test.html"}`,
			`{"file":"/tmp/test.pdf","status":"success"}`,
			http.StatusOK,
		},
		{
			`{"options":"--title Test","type":"string","string":"hello world!"}`,
			`{"file":"/tmp/output.pdf","status":"success"}`,
			http.StatusOK,
		},
		{
			`{"options":"--title Test","type":"url","url":"https://www.netzwerkorange.de/en"}`,
			`{"file":"/tmp/www.netzwerkorange.de.pdf","status":"success"}`,
			http.StatusOK,
		},
		{
			`{"options":"--title Test","type":"file","file":"not_here"}`,
			`{"error":"Input file not found. The file must reside in the shared dir."}`,
			http.StatusInternalServerError,
		},
		{
			`{"options":"--title Test","type":"file"}`,
			`{"error":"No filename provided."}`,
			http.StatusBadRequest,
		},
		{
			`{"options":"--title Test","type":"string","string":""}`,
			`{"error":"No string provided."}`,
			http.StatusBadRequest,
		},
		{
			`{"options":"--title Test","type":"file","string":"something"}`,
			`{"error":"No filename provided."}`,
			http.StatusBadRequest,
		},
		{
			`{"type":"unavailable"}`,
			`{"error":"Type not available."}`,
			http.StatusBadRequest,
		},
		{
			`{}`,
			`{"error":"Type not available."}`,
			http.StatusBadRequest,
		},
		{
			`{"type":"url", "url":"inval.id"}`,
			`{"error":"Error fetching URL."}`,
			http.StatusInternalServerError,
		},
		{
			`{"type":"url"}`,
			`{"error":"No URL provided."}`,
			http.StatusBadRequest,
		},
		{
			`{invalid_json}`,
			`{"error":"Invalid payload."}`,
			http.StatusBadRequest,
		},
	}

	ioutil.WriteFile("/tmp/test.html", []byte("Hello world!"), 0644)

	for _, tt := range tests {
		rr := httptest.NewRecorder()
		r := bytes.NewReader([]byte(tt.payload))

		req, err := http.NewRequest("POST", "/", r)
		if err != nil {
			t.Fatal(err)
		}

		handler := http.HandlerFunc(s.handleRequest)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != tt.status {
			t.Errorf("unexpected status code:\nhave %v\nwant %v\n", status, tt.status)
		}

		if rr.Body.String() != tt.expected {
			t.Errorf("unexpected body:\nhave %v\nwant %v\n", rr.Body.String(), tt.expected)
		}
	}
}
