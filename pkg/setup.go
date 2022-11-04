package pkg

import (
	"github.com/MarcusGoldschmidt/scadagobr/pkg/events"
	customLogger "github.com/MarcusGoldschmidt/scadagobr/pkg/logger"
	gorm2 "github.com/MarcusGoldschmidt/scadagobr/pkg/persistence/gorm"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/providers"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/purge"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/runtime"
	scadaServer "github.com/MarcusGoldschmidt/scadagobr/pkg/server"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"time"
)

func DefaultScadagobr(opt *ScadagobrOptions) (*Scadagobr, error) {
	loggerImp := customLogger.NewSimpleLogger("GOBR", os.Stdout)
	hubManager := events.NewHubManagerImpl(loggerImp)

	db, err := gorm.Open(postgres.Open(opt.PostgresConnectionString), &gorm.Config{
		Logger: customLogger.NewGormLogger(),
	})
	if err != nil {
		return nil, err
	}

	// Providers
	timeProvider := providers.TimeProviderFromTimeZone(opt.Location)

	// Persistence
	dataSourcePersistence := gorm2.NewDataSourcePersistenceGormImpl(db)
	dataPointPersistence := gorm2.NewDataPointPersistenceGormImpl(db, hubManager)
	userPersistence := gorm2.NewUserPersistenceImp(db)
	viewPersistence := gorm2.NewViewPersistenceGormImpl(db)

	// Runtime manager
	runtimeManager := runtime.NewRuntimeManager(loggerImp, dataPointPersistence, hubManager, timeProvider)
	purgeManager := purge.NewManager(dataPointPersistence, dataSourcePersistence, timeProvider, loggerImp, time.Hour)

	scada := &Scadagobr{
		RuntimeManager:        runtimeManager,
		Logger:                loggerImp,
		Db:                    db,
		Option:                opt,
		router:                mux.NewRouter(),
		JwtHandler:            SetupJwtHandler(opt, userPersistence),
		userPersistence:       userPersistence,
		dataSourcePersistence: dataSourcePersistence,
		dataPointPersistence:  dataPointPersistence,
		viewPersistence:       viewPersistence,
		internalRoute:         scadaServer.NewRouter(),
		purgeManager:          purgeManager,
		HubManager:            hubManager,
		timeProvider:          timeProvider,
	}

	scada.setRouters()

	err = scada.setServer()
	if err != nil {
		return nil, err
	}

	return scada, nil
}
