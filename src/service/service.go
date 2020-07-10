package service

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type (
	Service struct {
		hostname string
		shared   string
		binary   string
		request  payload
	}

	payload struct {
		Options string `json:"options"`
		Type    string `json:"type"`
		File    string `json:"file"`
		String  string `json:"string"`
		URL     string `json:"url"`
	}

	Map map[string]string
)

func New(hostname string) *Service {
	s := &Service{
		hostname: hostname,
		shared:   "/app/shared",
		binary:   "/usr/local/bin/wkhtmltopdf",
	}

	if _, err := os.Stat(s.shared); os.IsNotExist(err) {
		log.Fatalf("Shared dir not found. Mount a volume to %s first. Error: %s",
			s.shared,
			err.Error(),
		)
	}

	return s
}

func (s Service) router() *mux.Router {
	r := mux.NewRouter()
	r.Use(s.jsonMW, s.recoverMW)
	r.HandleFunc("/", s.handleRequest).Methods(http.MethodPost)
	return r
}

func (s Service) Serve() {
	srv := http.Server{
		Addr:         s.hostname,
		Handler:      s.router(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Info("PDF service listening at", s.hostname)
	log.Fatalln(srv.ListenAndServe())
}
