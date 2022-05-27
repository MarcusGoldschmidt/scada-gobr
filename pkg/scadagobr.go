package pkg

import (
	"context"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"scadagobr/pkg/auth"
	"scadagobr/pkg/logger"
	"scadagobr/pkg/models"
	"scadagobr/pkg/runtime"
	scadaServer "scadagobr/pkg/server"
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

	go func() {
		err := s.server.ListenAndServe()
		if err != nil {
			s.Logger.Errorf("%s", err.Error())
		}
	}()

	return nil
}

func (s *Scadagobr) Shutdown(ctx context.Context) {
	err := s.server.Shutdown(ctx)
	if err != nil {
		return
	}

	s.RuntimeManager.StopAll(ctx)
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
