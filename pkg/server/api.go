package server

import (
	"context"
	"embed"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/auth"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/logger"
	"github.com/gorilla/mux"
	"io"
	"io/fs"
	"net/http"
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

		r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			path := r.URL.Path

			if path[0] == '/' {
				path = path[1:]
			}

			_, err := files.Open(path)

			if err != nil {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.WriteHeader(http.StatusOK)

				file, _ := files.Open("index.html")

				response, err := io.ReadAll(file)

				_, err = w.Write(response)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
				return
			}

			http.FileServer(http.FS(files)).ServeHTTP(w, r)
		})
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
