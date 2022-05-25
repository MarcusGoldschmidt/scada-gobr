package server

import (
	"embed"
	"github.com/gorilla/mux"
	"io/fs"
	"net/http"
)

//go:generate cp -r ../../scadagobr-client/public ./public
//go:embed public
var spa embed.FS

func SetupRouters(r *mux.Router, devMode bool) error {
	if !devMode {
		files, err := fs.Sub(spa, "public")
		if err != nil {
			return err
		}

		r.Handle("/", http.FileServer(http.FS(files)))
	}

	return nil
}
