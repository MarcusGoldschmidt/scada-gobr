package pkg

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/buffers"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/events"
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
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func DefaultScadagobr(opt *ScadagobrOptions) (*Scadagobr, error) {
	ctx := context.Background()

	bufferSize := buffers.NewMaxBucketBuffer(32, buffers.MB*10)

	logOutput := io.MultiWriter(os.Stderr, bufferSize)

	loggerImp := custonLogger.NewSimpleLogger("GOBR", logOutput)
	hubManager := events.NewHubManagerImpl(loggerImp)

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
	dataSourcePersistence := gorm2.NewDataSourcePersistenceGormImpl(db)
	dataPointPersistence := gorm2.NewDataPointPersistenceGormImpl(db, hubManager)
	userPersistence := gorm2.NewUserPersistenceImp(db)
	viewPersistence := gorm2.NewViewPersistenceGormImpl(db)

	scadaRouter := scadaServer.NewRouter()

	// Runtime manager
	runtimeManager := runtime.NewRuntimeManager(loggerImp, dataPointPersistence)
	runtimeManager.WithTimeProvider(providers.UtcTimeProvider{})

	// Providers
	timeProvider := providers.UtcTimeProvider{}

	// Route to net server datasource
	r := mux.NewRouter()
	r.Handle("/api/datasource/integration", scadaRouter)

	jwtHandler := SetupJwtHandler(opt, userPersistence)

	purgeManager := purge.NewManager(dataPointPersistence, dataSourcePersistence, providers.UtcTimeProvider{}, loggerImp, time.Hour)

	scada := &Scadagobr{
		RuntimeManager:        runtimeManager,
		Logger:                loggerImp,
		Db:                    db,
		Option:                opt,
		router:                r,
		userPersistence:       userPersistence,
		JwtHandler:            jwtHandler,
		dataSourcePersistence: dataSourcePersistence,
		dataPointPersistence:  dataPointPersistence,
		internalRoute:         scadaRouter,
		purgeManager:          purgeManager,
		HubManager:            hubManager,
		viewPersistence:       viewPersistence,
		inMemoryLogs:          bufferSize,
		timeProvider:          timeProvider,
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

	scada.RuntimeManager.AddDataSourceManager(datasourceManagers...)

	return scada, nil
}
