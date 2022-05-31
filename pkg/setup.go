package pkg

import (
	"github.com/MarcusGoldschmidt/scadagobr/pkg/logger"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/models"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/persistence"
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

	userPersistence := persistence.NewUserPersistenceImp(db)

	jwtHandler := SetupJwtHandler(opt, userPersistence)

	scada := &Scadagobr{
		RuntimeManager:  runtimeManager,
		Logger:          simpleLog,
		Db:              db,
		Option:          opt,
		router:          r,
		UserPersistence: userPersistence,
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
