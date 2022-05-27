package pkg

import (
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"os"
	"scadagobr/pkg/logger"
	"scadagobr/pkg/models"
	"scadagobr/pkg/persistence"
	"scadagobr/pkg/providers"
	"scadagobr/pkg/runtime"
	scadaServer "scadagobr/pkg/server"
	"strconv"
	"time"
)

func DefaultScadagobr(opt *ScadagobrOptions) (*Scadagobr, error) {
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

	datasource := LoadDataSourceRuntimeManager(db, scadaRouter)

	runtimeManager.AddDataSource(datasource...)

	r := mux.NewRouter()

	r.Handle("/api/datasource/integration", scadaRouter)

	simpleLog := logger.NewSimpleLogger("SCADA", os.Stdout)

	scada := &Scadagobr{
		RuntimeManager: runtimeManager,
		Logger:         simpleLog,
		Db:             db,
		Option:         opt,
		router:         r,
	}

	scada.setRouters()

	err = scadaServer.SetupSpa(scada.router, opt.DevMode)
	if err != nil {
		return nil, err
	}

	scada.server = &http.Server{
		Handler:      scada.router,
		Addr:         opt.Address + ":" + strconv.Itoa(opt.Port),
		TLSConfig:    opt.TLSConfig,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	return scada, nil
}
