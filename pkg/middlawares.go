package pkg

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (s *Scadagobr) setupCors() {
	s.router.Use(mux.CORSMethodMiddleware(s.router))

	s.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
			if r.Method == "OPTIONS" {
				return
			}
			next.ServeHTTP(w, r)
		})
	})
}

func (s *Scadagobr) setupLogs() {
	s.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s.Logger.Infof("request at %s", r.RequestURI)
			next.ServeHTTP(w, r)
		})
	})
}
