package pkg

import (
	"github.com/MarcusGoldschmidt/scadagobr/pkg/events"
	customLogger "github.com/MarcusGoldschmidt/scadagobr/pkg/logger"
	gorm2 "github.com/MarcusGoldschmidt/scadagobr/pkg/persistence/gorm"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/providers"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/purge"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/runtime"
	scadaServer "github.com/MarcusGoldschmidt/scadagobr/pkg/server"
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
		JwtHandler:            SetupJwtHandler(opt, userPersistence),
		UserPersistence:       userPersistence,
		DataSourcePersistence: dataSourcePersistence,
		DataPointPersistence:  dataPointPersistence,
		ViewPersistence:       viewPersistence,
		InternalRouter:        scadaServer.NewRouter(),
		PurgeManager:          purgeManager,
		HubManager:            hubManager,
		TimeProvider:          timeProvider,
	}

	return scada, nil
}
