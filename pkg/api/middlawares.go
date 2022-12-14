package api

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/auth"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/providers"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/util"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/attribute"
	"net/http"
	"strings"
)

func (s *ScadaApi) setupCors() {
	s.router.Use(mux.CORSMethodMiddleware(s.router))

	s.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
			if r.Method == "OPTIONS" {
				w.WriteHeader(200)
				return
			}
			next.ServeHTTP(w, r)
		})
	})
}

func (s *ScadaApi) setupProviders() {
	s.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), providers.TimeProviderCtxKey, s.TimeProvider)
			r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	})
}

func (s *ScadaApi) setupLogs() {
	s.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s.Logger.Tracef("Request at %s", r.RequestURI)
			ctx, span := util.Tracer.Start(r.Context(), r.Method)
			defer span.End()

			span.SetAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.host", r.Host),
			)

			r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	})
}

func (s *ScadaApi) jwtMiddleware(function RequestHandlerFunction) RequestHandlerFunction {
	return func(scada *ScadaApi, w http.ResponseWriter, r *http.Request) {
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

func (s *ScadaApi) adminMiddleware(function RequestHandlerFunction) RequestHandlerFunction {
	return func(scada *ScadaApi, w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		claims, err := auth.GetUserFromContext(ctx)
		if err != nil {
			s.respondError(r.Context(), w, err)
			return
		}

		if !claims.Admin {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		function(scada, w, r)
	}
}

func (s *ScadaApi) authAndIsAdminMiddleware(function RequestHandlerFunction) RequestHandlerFunction {
	return MultipleMiddleware(function, s.jwtMiddleware, s.adminMiddleware)
}

func MultipleMiddleware(h RequestHandlerFunction, m ...func(RequestHandlerFunction) RequestHandlerFunction) RequestHandlerFunction {
	if len(m) < 1 {
		return h
	}

	wrapped := h

	for i := len(m) - 1; i >= 0; i-- {
		wrapped = m[i](wrapped)
	}

	return wrapped
}
