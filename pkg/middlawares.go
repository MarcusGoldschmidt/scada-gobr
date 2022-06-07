package pkg

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/auth"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
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

func (s *Scadagobr) jwtMiddleware(function RequestHandlerFunction) RequestHandlerFunction {
	return func(scada *Scadagobr, w http.ResponseWriter, r *http.Request) {
		reqToken := r.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")
		reqToken = splitToken[1]

		claims, err := scada.JwtHandler.ValidateJwt(reqToken)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), auth.ContextClaimsKey, claims)

		r = r.WithContext(ctx)

		function(scada, w, r)
	}
}

func (s *Scadagobr) adminMiddleware(function RequestHandlerFunction) RequestHandlerFunction {
	return func(scada *Scadagobr, w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		claims, err := auth.GetUserFromContext(ctx)
		if err != nil {
			s.respondError(w, err)
			return
		}

		if !claims.Admin {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		function(scada, w, r)
	}
}
