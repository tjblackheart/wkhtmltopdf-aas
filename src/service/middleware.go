package service

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (s Service) jsonMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (s Service) recoverMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Error(err)
				w.Header().Set("Connection", "close")
				return
			}
		}()

		next.ServeHTTP(w, r)
	})
}
