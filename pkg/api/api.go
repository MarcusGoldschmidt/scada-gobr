package api

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/logger"
	"github.com/gorilla/mux"
	"net/http"
)

type ScadaApi struct {
	*pkg.Scadagobr
	server *http.Server
	router *mux.Router
	Logger logger.Logger
}

func DefaultScadaApi(opt *pkg.ScadagobrOptions) (*ScadaApi, error) {
	scada, err := pkg.DefaultScadagobr(opt)
	if err != nil {
		return nil, err
	}

	return &ScadaApi{
		Scadagobr: scada,
		server:    nil,
		router:    mux.NewRouter(),
		Logger:    scada.Logger,
	}, nil
}

func (s *ScadaApi) SetupAndRun(ctx context.Context) error {
	err := s.Scadagobr.SetupAndRun(ctx)
	if err != nil {
		return err
	}

	err = s.setServer()
	if err != nil {
		return err
	}

	return s.ListenAndServeHttp(ctx)
}

func (s *ScadaApi) ListenAndServeHttp(ctx context.Context) error {
	protocol := "https://"
	if s.server.TLSConfig == nil {
		protocol = "http://"
	}

	s.Logger.Infof("Start HTTP server with address: %s%s", protocol, s.server.Addr)

	err := s.server.ListenAndServe()
	if err != nil {
		s.Logger.Infof("%s", err.Error())
		return err
	}
	return nil
}

func (s *ScadaApi) Shutdown(ctx context.Context) {
	err := s.server.Shutdown(ctx)
	if err != nil {
		return
	}

	s.Scadagobr.Shutdown(ctx)
}
