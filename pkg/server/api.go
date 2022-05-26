package server

import (
	"context"
	"embed"
	"github.com/gorilla/mux"
	"io/fs"
	"net/http"
	"scadagobr/pkg/auth"
	"scadagobr/pkg/logger"
	"strings"
)

//go:generate cp -r ../../scadagobr-client/public ./public
//go:embed public
var spa embed.FS

func SetupSpa(r *mux.Router, devMode bool) error {
	if !devMode {
		files, err := fs.Sub(spa, "public")
		if err != nil {
			return err
		}

		r.PathPrefix("/").Handler(http.FileServer(http.FS(files)))
	}

	return nil
}

type JwtHttpFunc struct {
	callback http.Handler
	handler  *auth.JwtHandler
	logger   logger.Logger
}

func NewJwtHttpFunc(callback http.Handler, handler *auth.JwtHandler, logger logger.Logger) *JwtHttpFunc {
	return &JwtHttpFunc{callback: callback, handler: handler, logger: logger}
}

func (j JwtHttpFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]

	claims, err := j.handler.ValidateJwt(reqToken)
	if err != nil {
		if err == auth.ErrUnauthorized {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		j.logger.Errorf("%s", err.Error())
		return
	}

	ctx := context.WithValue(r.Context(), auth.ContexClaimsKey, claims)

	r = r.WithContext(ctx)

	j.callback.ServeHTTP(w, r)
}
