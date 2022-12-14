package api

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/server"
	"net/http"
)

func HealthCheckHandler(scada *ScadaApi, w http.ResponseWriter, r *http.Request) {
	server.Handler(
		server.WithChecker(
			"database", server.CheckerFunc(func(ctx context.Context) error {
				sqlDB, err := scada.Db.DB()

				if err != nil {
					return err
				}

				err = sqlDB.Ping()
				if err != nil {
					return err
				}

				return nil
			}),
		),
	).ServeHTTP(w, r)
}
