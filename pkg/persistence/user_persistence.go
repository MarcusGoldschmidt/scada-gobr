package persistence

import (
	"github.com/google/uuid"
	"scadagobr/pkg/models"
)

type UserPersistence interface {
	GetUser(uuid.UUID) (*models.User, error)
}
