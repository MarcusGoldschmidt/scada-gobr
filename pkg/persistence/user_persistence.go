package persistence

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"scadagobr/pkg/models"
)

type UserPersistence interface {
	GetUser(uuid.UUID) (*models.User, error)
}

type UserPersistenceImp struct {
	db *gorm.DB
}

func (u UserPersistenceImp) GetUser(u2 uuid.UUID) (*models.User, error) {
	var user *models.User
	result := u.db.Model(&models.User{}).Where(models.User{Name: "admin"}).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func NewUserPersistenceImp(db *gorm.DB) *UserPersistenceImp {
	return &UserPersistenceImp{db: db}
}
