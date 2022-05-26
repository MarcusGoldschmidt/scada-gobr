package pkg

import (
	"context"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"os"
	"scadagobr/pkg/auth"
	"scadagobr/pkg/datasources"
	"scadagobr/pkg/logger"
	"scadagobr/pkg/models"
	"scadagobr/pkg/persistence"
	"scadagobr/pkg/providers"
	"scadagobr/pkg/runtime"
	scadaServer "scadagobr/pkg/server"
	"strconv"
	"time"
)

type Scadagobr struct {
	RuntimeManager *runtime.RuntimeManager
	Logger         logger.Logger
	Db             *gorm.DB
	Option         *ScadagobrOptions

	server *http.Server
	router *mux.Router
}

func (s *Scadagobr) Setup() error {

	var user *models.User
	result := s.Db.Model(&models.User{}).Where(models.User{Name: "admin"}).First(&user)

	if result.RowsAffected == 0 {
		result := s.Db.Create(&models.User{
			ID:            uuid.New(),
			Name:          "admin",
			Administrator: true,
			Email:         &s.Option.AdminEmail,
			HomeUrl:       "/",
			PasswordHash:  auth.MakeHash(s.Option.AdminPassword),
		})

		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

func (s *Scadagobr) Run(ctx context.Context) error {
	s.RuntimeManager.RunAll(ctx)

	s.Logger.Infof("Start HTTP server with address: %s", s.server.Addr)
	err := s.server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

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

	datasource := datasources.LoadDataSources(db, scadaRouter)

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

type RequestHandlerFunction func(scada *Scadagobr, w http.ResponseWriter, r *http.Request)

func (s *Scadagobr) handleRequest(function RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		function(s, w, r)
	}
}

func (s *Scadagobr) handleJwtRequest(jwtHandler *auth.JwtHandler, function RequestHandlerFunction) http.Handler {

	callback := s.handleRequest(function)

	httpFunction := scadaServer.NewJwtHttpFunc(callback, jwtHandler, s.Logger)

	return httpFunction
}

func (s *Scadagobr) setRouters() {
	//userPersistence := persistence.NewUserPersistenceImp(s.Db)

	//jwtHandler := SetupJwtHandler(s.Option, userPersistence)

	s.get("/api/healthcheck", s.handleRequest(Func))
}

func (s *Scadagobr) get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	s.router.HandleFunc(path, f).Methods("GET")
}

// Post wraps the router for POST method
func (s *Scadagobr) post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	s.router.HandleFunc(path, f).Methods("POST")
}

// Put wraps the router for PUT method
func (s *Scadagobr) put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	s.router.HandleFunc(path, f).Methods("PUT")
}

// Delete wraps the router for DELETE method
func (s *Scadagobr) delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	s.router.HandleFunc(path, f).Methods("DELETE")
}
