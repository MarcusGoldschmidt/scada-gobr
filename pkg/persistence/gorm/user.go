package gorm

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/auth"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/models"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserPersistenceGormImpl struct {
	db *gorm.DB
}

func (u UserPersistenceGormImpl) CreateAdminUser(ctx context.Context, email, password string) error {
	db := u.db.WithContext(ctx)

	var user *models.User
	result := db.Model(&models.User{}).Where(models.User{Name: "admin"}).First(&user)

	hash, _ := auth.MakeHash(password)

	if result.RowsAffected == 0 {
		result := db.Create(&models.User{
			ID:            uuid.New(),
			Name:          "admin",
			Administrator: true,
			Email:         &email,
			HomeUrl:       "/",
			PasswordHash:  hash,
		})

		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

func NewUserPersistenceGormImpl(db *gorm.DB) *UserPersistenceGormImpl {
	return &UserPersistenceGormImpl{db: db}
}

func (u UserPersistenceGormImpl) GetUsers(ctx context.Context, request *shared.PaginationRequest) ([]*models.User, error) {
	db := u.db.WithContext(ctx)

	return listPaginate[models.User](db, request)
}

func (u UserPersistenceGormImpl) IsValidUsernameForUser(ctx context.Context, name string, id uuid.UUID) (bool, error) {
	db := u.db.WithContext(ctx)

	var count int64
	err := db.Model(&models.User{}).Where("id <> ? AND name = ?", id.String(), name).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func NewUserPersistenceImp(db *gorm.DB) *UserPersistenceGormImpl {
	return &UserPersistenceGormImpl{db: db}
}

func (u UserPersistenceGormImpl) CreateUser(ctx context.Context, user *models.User) error {
	db := u.db.WithContext(ctx)
	return db.Create(user).Error
}

func (u UserPersistenceGormImpl) UpdateUser(ctx context.Context, user *models.User) error {
	db := u.db.WithContext(ctx)
	return db.Save(user).Error
}

func (u UserPersistenceGormImpl) DeleteUser(ctx context.Context, id uuid.UUID) error {
	db := u.db.WithContext(ctx)
	return db.Delete(&models.User{ID: id}).Error
}

func (u UserPersistenceGormImpl) GetUserByName(ctx context.Context, name string) (*models.User, error) {
	db := u.db.WithContext(ctx)

	var user *models.User
	result := db.Model(&models.User{}).Where(models.User{Name: name}).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (u UserPersistenceGormImpl) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	db := u.db.WithContext(ctx)

	var user *models.User
	result := db.Model(&models.User{}).Where(models.User{Email: &email}).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (u UserPersistenceGormImpl) GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	db := u.db.WithContext(ctx)

	var user *models.User
	result := db.Model(&models.User{}).Where(models.User{ID: id}).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}
