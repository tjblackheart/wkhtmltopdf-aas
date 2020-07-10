package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

func (s Service) handleRequest(w http.ResponseWriter, r *http.Request) {
	if err := json.NewDecoder(r.Body).Decode(&s.request); err != nil {
		log.Error(err)
		s.response(w, http.StatusBadRequest, Map{"error": "Invalid payload."})
		return
	}

	switch strings.ToLower(s.request.Type) {
	case "string":
		s.processString(w, r)
	case "url":
		s.processURL(w, r)
	case "file":
		s.processFile(w, r)
	default:
		s.response(w, http.StatusBadRequest, Map{"error": "Type not available."})
	}
}

//

func (s Service) processFile(w http.ResponseWriter, r *http.Request) {
	if s.request.File == "" {
		s.response(w, http.StatusBadRequest, Map{"error": "No filename provided."})
		return
	}

	file := filepath.Base(s.request.File)
	input := fmt.Sprintf("%s/%s", s.shared, file)
	if _, err := os.Stat(input); os.IsNotExist(err) {
		log.Error(err)
		s.response(w, http.StatusInternalServerError, Map{
			"error": "Input file not found. The file must reside in the shared dir.",
		})
		return
	}

	target := strings.TrimSuffix(file, filepath.Ext(file)) + ".pdf"
	output := fmt.Sprintf("%s/%s", s.shared, target)

	if s.request.Options == "" {
		s.request.Options = "-q"
	}

	opts := strings.Split(s.request.Options, " ")
	opts = append(opts, input)
	opts = append(opts, output)

	out, err := exec.Command(s.binary, opts...).CombinedOutput()
	if err != nil {
		log.Error(string(out), err)
		s.response(w, http.StatusInternalServerError, Map{"error": err.Error()})
		return
	}

	s.response(w, http.StatusOK, Map{
		"status": "success",
		"file":   output,
	})
}

func (s Service) processString(w http.ResponseWriter, r *http.Request) {
	if s.request.String == "" {
		s.response(w, http.StatusBadRequest, Map{"error": "No string provided."})
		return
	}

	u := &url.URL{}

	// from processURL
	if s.request.Type == "url" {
		u, _ = url.Parse(s.request.URL)
	}

	output := "output.html"

	if u.Host != "" {
		output = fmt.Sprintf("%s.html", u.Host)

		abs := fmt.Sprintf("%s://%s", u.Scheme, u.Host)
		replacements := map[string]string{
			"src=\"/":   "src=\"%s/",
			"src=\"//":  "src=\"%s/",
			"href=\"/":  "href=\"%s/",
			"href=\"//": "href=\"%s/",
		}

		for search, rpl := range replacements {
			s.request.String = strings.ReplaceAll(s.request.String, search, fmt.Sprintf(rpl, abs))
		}
	}

	s.request.File = fmt.Sprintf("%s/%s", s.shared, output)
	if err := ioutil.WriteFile(s.request.File, []byte(s.request.String), 0644); err != nil {
		log.Error(err)
		s.response(w, http.StatusInternalServerError, Map{"error": "Error writing tmp file."})
		return
	}

	s.processFile(w, r)
}

func (s Service) processURL(w http.ResponseWriter, r *http.Request) {
	if s.request.URL == "" {
		s.response(w, http.StatusBadRequest, Map{"error": "No URL provided."})
		return
	}

	resp, err := http.Get(s.request.URL)
	if err != nil {
		log.Error(err)
		s.response(w, http.StatusInternalServerError, Map{"error": "Error fetching URL."})
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}

	s.request.String = string(body)
	s.processString(w, r)
}

//

func (s Service) response(w http.ResponseWriter, status int, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		log.Error(err)
	}

	w.WriteHeader(status)
	w.Write(b)
}
