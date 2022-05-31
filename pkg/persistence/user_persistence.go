package persistence

import (
	"github.com/MarcusGoldschmidt/scadagobr/pkg/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserPersistence interface {
	GetUserById(uuid.UUID) (*models.User, error)
	GetUserByEmail(string) (*models.User, error)
	GetUserByName(string) (*models.User, error)
}

type UserPersistenceImp struct {
	db *gorm.DB
}

func (u UserPersistenceImp) GetUserByName(name string) (*models.User, error) {
	var user *models.User
	result := u.db.Model(&models.User{}).Where(models.User{Name: name}).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (u UserPersistenceImp) GetUserByEmail(email string) (*models.User, error) {
	var user *models.User
	result := u.db.Model(&models.User{}).Where(models.User{Email: &email}).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (u UserPersistenceImp) GetUserById(id uuid.UUID) (*models.User, error) {
	var user *models.User
	result := u.db.Model(&models.User{}).Where(models.User{ID: id}).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func NewUserPersistenceImp(db *gorm.DB) *UserPersistenceImp {
	return &UserPersistenceImp{db: db}
}
