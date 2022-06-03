package pkg

import (
	"context"
	custonLogger "github.com/MarcusGoldschmidt/scadagobr/pkg/logger"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/models"
	gorm2 "github.com/MarcusGoldschmidt/scadagobr/pkg/persistence/gorm"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/providers"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/purge"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/runtime"
	scadaServer "github.com/MarcusGoldschmidt/scadagobr/pkg/server"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strconv"
	"time"
)

func DefaultScadagobr(opt *ScadagobrOptions) (*Scadagobr, error) {
	ctx := context.Background()

	loggerImp := custonLogger.NewSimpleLogger("RUNTIME-MANAGER", os.Stdout)

	db, err := gorm.Open(postgres.Open(opt.PostgresConnectionString), &gorm.Config{
		Logger: custonLogger.NewGormLogger(),
	})
	if err != nil {
		return nil, err
	}

	err = models.AutoMigration(db)
	if err != nil {
		return nil, err
	}

	// Persistence
	persistenceImp := gorm2.NewDataPointPersistenceGormImpl(db)
	dataSourcePersistence := gorm2.NewDataSourcePersistenceGormImpl(db)
	dataPointPersistence := gorm2.NewDataPointPersistenceGormImpl(db)
	userPersistence := gorm2.NewUserPersistenceImp(db)

	scadaRouter := scadaServer.NewRouter()

	// Runtime manager
	runtimeManager := runtime.NewRuntimeManager(loggerImp, persistenceImp)
	runtimeManager.WithTimeProvider(providers.UtcTimeProvider{})

	// Route to http server datasource
	r := mux.NewRouter()
	r.Handle("/api/datasource/integration", scadaRouter)

	jwtHandler := SetupJwtHandler(opt, userPersistence)
	simpleLog := custonLogger.NewSimpleLogger("SCADA", os.Stdout)

	purgeManager := purge.NewManager(dataPointPersistence, dataSourcePersistence, providers.UtcTimeProvider{}, loggerImp, time.Hour)

	scada := &Scadagobr{
		RuntimeManager:        runtimeManager,
		Logger:                simpleLog,
		Db:                    db,
		Option:                opt,
		router:                r,
		userPersistence:       userPersistence,
		JwtHandler:            jwtHandler,
		dataSourcePersistence: dataSourcePersistence,
		dataPointPersistence:  dataPointPersistence,
		internalRoute:         scadaRouter,
		purgeManager:          purgeManager,
	}

	scada.setRouters()

	if !opt.DevMode {
		err = scadaServer.SetupSpa(scada.router)
		if err != nil {
			return nil, err
		}
	}

	scada.server = &http.Server{
		Handler:      scada.router,
		Addr:         opt.Address + ":" + strconv.Itoa(opt.Port),
		TLSConfig:    opt.TLSConfig,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	datasourceManagers, err := LoadDataSourcesRuntimeManager(ctx, scada)
	if err != nil {
		return nil, err
	}

	scada.RuntimeManager.AddDataSource(datasourceManagers...)

	return scada, nil
}
