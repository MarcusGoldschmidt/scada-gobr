package pkg

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/auth"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/events"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/logger"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/models"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/persistence"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/providers"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/purge"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/runtime"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/server"
	"gorm.io/gorm"
)

type Scadagobr struct {
	RuntimeManager *runtime.Manager
	Logger         logger.Logger
	Db             *gorm.DB
	Option         *ScadagobrOptions

	// Persistence
	UserPersistence       persistence.UserPersistence
	DataPointPersistence  persistence.DataPointPersistence
	DataSourcePersistence persistence.DataSourcePersistence
	ViewPersistence       persistence.ViewPersistence

	JwtHandler *auth.JwtHandler

	TimeProvider providers.TimeProvider

	InternalRouter *server.Router

	PurgeManager *purge.Manager
	HubManager   events.HubManager

	// Created after the server is started
	shutdownContext func()
}

func (s *Scadagobr) Setup(ctx context.Context) error {
	s.Logger.Infof("VERSION: %s", App)
	s.Logger.Infof("VERSION: %s", Version)
	s.Logger.Infof("COMMIT: %s", Commit)
	s.Logger.Infof("BUILT AT: %s", BuiltAt)

	ctx, trace := s.Trace(ctx, "Setting up Scadagobr")
	defer trace.End()

	err := models.AutoMigration(s.Db)
	if err != nil {
		return err
	}

	err = s.UserPersistence.CreateAdminUser(ctx, s.Option.AdminEmail, s.Option.AdminPassword)
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
	ctx = context.WithValue(ctx, providers.TimeProviderCtxKey, s.TimeProvider)
	ctx, s.shutdownContext = context.WithCancel(ctx)

	s.Logger.Infof("Starting Scadagobr")

	go s.PurgeManager.Work(ctx)

	err := s.RuntimeManager.RunAll(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Scadagobr) SetupAndRun(ctx context.Context) error {
	err := s.Setup(ctx)

	if err != nil {
		return err
	}

	err = s.Run(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Scadagobr) Shutdown(ctx context.Context) {
	s.shutdownContext()

	s.RuntimeManager.StopAll(ctx)
	s.HubManager.ShutDown(ctx)
}
