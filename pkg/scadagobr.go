package pkg

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/auth"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/logger"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/models"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/persistence"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/runtime"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/server"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

type Scadagobr struct {
	RuntimeManager *runtime.RuntimeManager
	Logger         logger.Logger
	Db             *gorm.DB
	Option         *ScadagobrOptions

	// Persistence
	userPersistence       persistence.UserPersistence
	dataPointPersistence  persistence.DataPointPersistence
	dataSourcePersistence persistence.DataSourcePersistence

	JwtHandler *auth.JwtHandler

	server *http.Server
	router *mux.Router

	internalRoute *server.Router
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
