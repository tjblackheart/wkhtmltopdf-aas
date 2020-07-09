package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type (
	service struct {
		hostname  string
		sharedDir string
		binary    string
		payload   struct {
			Options string `json:"options"`
			Type    string `json:"type"`
			File    string `json:"file"`
			String  string `json:"string"`
			URL     string `json:"url"`
		}
	}

	data map[string]string
)

func main() {
	s := &service{
		hostname:  os.Getenv("APP_HOST"),
		sharedDir: "/app/shared",
		binary:    "/usr/local/bin/wkhtmltopdf",
	}

	r := mux.NewRouter()
	r.HandleFunc("/", s.handleRequest).Methods(http.MethodPost)

	srv := http.Server{
		Addr:         s.hostname,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	log.Info("PDF service listening at", s.hostname)
	log.Fatalln(srv.ListenAndServe())
}

func (s service) handleRequest(w http.ResponseWriter, r *http.Request) {
	if err := json.NewDecoder(r.Body).Decode(&s.payload); err != nil {
		log.Error(err)
		s.response(w, http.StatusInternalServerError, data{"error": err.Error()})
		return
	}

	if _, err := os.Stat(s.sharedDir); os.IsNotExist(err) {
		log.Error(err)
		s.response(w, http.StatusInternalServerError, data{
			"error": "Shared dir does not exist. Mount a volume to /app/shared first.",
		})
		return
	}

	switch strings.ToLower(s.payload.Type) {
	case "string":
		s.processString(w, r)
	case "url":
		s.processURL(w, r)
	case "file":
		s.processFile(w, r)
	default:
		s.response(w, http.StatusBadRequest, data{"error": "Type not available."})
	}
}

//

func (s service) processFile(w http.ResponseWriter, r *http.Request) {
	file := filepath.Base(s.payload.File)
	input := fmt.Sprintf("%s/%s", s.sharedDir, file)
	if _, err := os.Stat(input); os.IsNotExist(err) {
		log.Error(err)
		s.response(w, http.StatusInternalServerError, data{
			"error": "Input file not found. The file must reside in the shared dir.",
		})
		return
	}

	// err := s.detectContentType(input)
	// if err != nil {
	// 	log.Error(err)
	// 	s.response(w, http.StatusInternalServerError, data{"error": "File not supported."})
	// 	return
	// }

	target := strings.TrimSuffix(file, filepath.Ext(file)) + ".pdf"
	output := fmt.Sprintf("%s/%s", s.sharedDir, target)

	opts := strings.Split(s.payload.Options, " ")
	opts = append(opts, input)
	opts = append(opts, output)

	out, err := exec.Command("/usr/local/bin/wkhtmltopdf", opts...).CombinedOutput()
	if err != nil {
		log.Error(string(out), err)
		s.response(w, http.StatusInternalServerError, data{"error": err.Error()})
		return
	}

	s.response(w, http.StatusOK, data{
		"status": "success",
		"file":   output,
	})
}

func (s service) processString(w http.ResponseWriter, r *http.Request) {
	if s.payload.String == "" {
		s.response(w, http.StatusBadRequest, data{"error": "No string provided."})
		return
	}

	var err error
	u := &url.URL{}

	if s.payload.Type == "url" && s.payload.URL != "" {
		u, err = url.Parse(s.payload.URL)
		if err != nil {
			log.Error(err)
			s.response(w, http.StatusBadRequest, data{"error": "Error parsing URL."})
			return
		}
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
			s.payload.String = strings.ReplaceAll(s.payload.String, search, fmt.Sprintf(rpl, abs))
		}
	}

	s.payload.File = fmt.Sprintf("%s/%s", s.sharedDir, output)
	if err := ioutil.WriteFile(s.payload.File, []byte(s.payload.String), 0644); err != nil {
		log.Error(err)
		s.response(w, http.StatusInternalServerError, data{"error": "Error writing tmp file."})
		return
	}

	s.processFile(w, r)
}

func (s service) processURL(w http.ResponseWriter, r *http.Request) {
	if s.payload.URL == "" {
		s.response(w, http.StatusBadRequest, data{"error": "No URL provided."})
		return
	}

	resp, err := http.Get(s.payload.URL)
	if err != nil {
		log.Error(err)
		s.response(w, http.StatusInternalServerError, data{"error": "Error fetching URL."})
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}

	s.payload.String = string(body)
	s.processString(w, r)
}

//

func (s service) response(w http.ResponseWriter, status int, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		log.Error(err)
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)
	w.Write(b)
}

func (s service) detectContentType(path string) error {
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
