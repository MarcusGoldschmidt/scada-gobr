package server

import (
	"embed"
	"github.com/gorilla/mux"
	"io"
	"io/fs"
	"net/http"
	"strings"
)

//go:generate cp -r ../../scadagobr-client/public ./public
//go:embed public
var spa embed.FS

func SetupSpa(r *mux.Router) error {
	files, err := fs.Sub(spa, "public")
	if err != nil {
		return err
	}

	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		path := r.URL.Path

		if path[0] == '/' {
			path = path[1:]
		}

		if strings.HasPrefix(path, "api") {
			w.WriteHeader(http.StatusNotFound)
			return
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

	return nil
}
