package api

import (
	scadaServer "github.com/MarcusGoldschmidt/scadagobr/pkg/server"
	"net/http"
	"strconv"
	"time"
)

func (s *ScadaApi) setServer() error {
	if s.server != nil {
		return nil
	}

	err := scadaServer.SetupSpa(s.router)
	if err != nil {
		return err
	}

	s.server = &http.Server{
		Handler:      s.router,
		Addr:         s.Option.Address + ":" + strconv.Itoa(s.Option.Port),
		TLSConfig:    s.Option.TLSConfig,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	return nil
}
