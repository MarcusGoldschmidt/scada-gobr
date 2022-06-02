package pkg

import (
	custonLogger "github.com/MarcusGoldschmidt/scadagobr/pkg/logger"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/models"
	gorm2 "github.com/MarcusGoldschmidt/scadagobr/pkg/persistence/gorm"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/persistence/in_memory"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/providers"
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
	loggerImp := custonLogger.NewSimpleLogger("RUNTIME-MANAGER", os.Stdout)
	persistenceImp := in_memory.NewInMemoryPersistence()
	runtimeManager := runtime.NewRuntimeManager(loggerImp, persistenceImp)

	runtimeManager.WithTimeProvider(providers.UtcTimeProvider{})

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

	scadaRouter := scadaServer.NewRouter()

	datasource := LoadDataSourceRuntimeManager(db, scadaRouter)

	runtimeManager.AddDataSource(datasource...)

	r := mux.NewRouter()

	r.Handle("/api/datasource/integration", scadaRouter)

	userPersistence := gorm2.NewUserPersistenceImp(db)

	jwtHandler := SetupJwtHandler(opt, userPersistence)

	simpleLog := custonLogger.NewSimpleLogger("SCADA", os.Stdout)

	scada := &Scadagobr{
		RuntimeManager:  runtimeManager,
		Logger:          simpleLog,
		Db:              db,
		Option:          opt,
		router:          r,
		userPersistence: userPersistence,
		JwtHandler:      jwtHandler,
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
