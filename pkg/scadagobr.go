package pkg

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/auth"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/events"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/logger"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/persistence"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/providers"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/purge"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/runtime"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/server"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

type Scadagobr struct {
	RuntimeManager *runtime.Manager
	Logger         logger.Logger
	Db             *gorm.DB
	Option         *ScadagobrOptions

	// Persistence
	userPersistence       persistence.UserPersistence
	dataPointPersistence  persistence.DataPointPersistence
	dataSourcePersistence persistence.DataSourcePersistence
	viewPersistence       persistence.ViewPersistence

	JwtHandler *auth.JwtHandler

	server *http.Server
	router *mux.Router

	timeProvider providers.TimeProvider

	internalRoute *server.Router

	purgeManager *purge.Manager
	HubManager   events.HubManager

	// Created after the server is started
	shutdownContext func()
}

func (s *Scadagobr) Setup(ctx context.Context) error {

	err := s.userPersistence.CreateAdminUser(ctx, s.Option.AdminEmail, s.Option.AdminPassword)
	if err != nil {
		return err
	}

	datasourceManagers, err := s.LoadDataSourcesRuntimeManager(ctx)
	if err != nil {
		return err
	}

	s.RuntimeManager.AddDataSourceManager(datasourceManagers...)

	return nil
}

func (s *Scadagobr) Run(ctx context.Context) error {
	ctx = context.WithValue(ctx, providers.TimeProviderCtxKey, s.timeProvider)
	ctx, s.shutdownContext = context.WithCancel(ctx)

	go s.purgeManager.Work(ctx)

	err := s.RuntimeManager.RunAll(ctx)
	if err != nil {
		return err
	}

	go s.ListenAndServeHttp(ctx)

	return nil
}

func (s *Scadagobr) ListenAndServeHttp(ctx context.Context) {
	protocol := "https://"
	if s.server.TLSConfig == nil {
		protocol = "http://"
	}

	s.Logger.Infof("Start HTTP server with address: %s%s", protocol, s.server.Addr)

	err := s.server.ListenAndServe()
	if err != nil {
		s.Logger.Errorf("%s", err.Error())
	}
}

func (s *Scadagobr) Shutdown(ctx context.Context) {
	s.shutdownContext()

	s.HubManager.ShutDown(ctx)

	err := s.server.Shutdown(ctx)
	if err != nil {
		return
	}

	s.RuntimeManager.StopAll(ctx)
}
