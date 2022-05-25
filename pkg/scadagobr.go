package pkg

import (
	"context"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"os"
	datasources2 "scadagobr/pkg/datasources"
	"scadagobr/pkg/logger"
	"scadagobr/pkg/models"
	"scadagobr/pkg/persistence"
	"scadagobr/pkg/providers"
	"scadagobr/pkg/runtime"
	scadaServer "scadagobr/pkg/server"
	"strconv"
	"time"
)

type ScadagobrServer interface {
	Run(ctx context.Context) error
}

type Scadagobr struct {
	rt     *runtime.RuntimeManager
	server *http.Server
	logger logger.Logger
}

func (s *Scadagobr) Run(ctx context.Context) error {
	s.rt.RunAll(ctx)

	s.logger.Infof("Start HTTP server with address: %s", s.server.Addr)
	err := s.server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func DefaultScadagobr(opt *ScadagobrOptions) (ScadagobrServer, error) {
	loggerImp := logger.NewSimpleLogger("RUNTIME-MANAGER", os.Stdout)
	persistenceImp := persistence.NewInMemoryPersistence()
	runtimeManager := runtime.NewRuntimeManager(loggerImp, persistenceImp)

	runtimeManager.WithTimeProvider(providers.UtcTimeProvider{})

	scadaRouter := scadaServer.NewRouter()

	db, err := gorm.Open(postgres.Open(opt.PostgresConnectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = models.AutoMigration(db)
	if err != nil {
		return nil, err
	}

	datasources := datasources2.LoadDataSources(db, scadaRouter)

	runtimeManager.AddDataSource(datasources...)

	r := mux.NewRouter()

	err = scadaServer.SetupRouters(r, opt.DevMode)
	if err != nil {
		return nil, err
	}

	r.Handle("/api/datasource/integration", scadaRouter)

	httpServer := &http.Server{
		Handler:   r,
		Addr:      opt.Address + ":" + strconv.Itoa(opt.Port),
		TLSConfig: opt.TLSConfig,

		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	simpleLog := logger.NewSimpleLogger("SCADA", os.Stdout)

	scada := ScadagobrServer(&Scadagobr{runtimeManager, httpServer, simpleLog})

	return scada, nil
}
