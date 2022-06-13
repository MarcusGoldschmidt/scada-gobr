package pkg

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/auth"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/buffers"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/events"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/logger"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/models"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/persistence"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/purge"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/runtime"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/server"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

type Scadagobr struct {
	RuntimeManager *runtime.Manager
	Logger         logger.Logger
	Db             *gorm.DB
	Option         *ScadagobrOptions

	// Persistence
	userPersistence       persistence.UserPersistence
	dataPointPersistence  persistence.DataPointPersistence
	dataSourcePersistence persistence.DataSourcePersistence
	viewPersistence       persistence.ViewPersistence

	JwtHandler *auth.JwtHandler

	server *http.Server
	router *mux.Router

	internalRoute *server.Router

	purgeManager *purge.Manager
	HubManager   events.HubManager
	wsManager    WsManager

	// Created after the server is started
	shutdownContext func()

	inMemoryLogs *buffers.MaxBucketBuffer
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
	ctx, s.shutdownContext = context.WithCancel(ctx)

	go s.purgeManager.Work(ctx)

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
	s.shutdownContext()

	s.HubManager.ShutDown(ctx)

	err := s.server.Shutdown(ctx)
	if err != nil {
		return
	}

	s.RuntimeManager.StopAll(ctx)
}
