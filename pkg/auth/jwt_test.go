package auth

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/models"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/providers"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"github.com/google/uuid"
	"testing"
	"time"
)

type fakeUserPersistence struct {
	mock *models.User
}

func (f fakeUserPersistence) GetUserById(ctx context.Context, u uuid.UUID) (*models.User, error) {
	return f.mock, nil
}

func (f fakeUserPersistence) GetUserByEmail(ctx context.Context, s string) (*models.User, error) {
	return f.mock, nil
}

func (f fakeUserPersistence) GetUserByName(ctx context.Context, s string) (*models.User, error) {
	return f.mock, nil
}

func (f fakeUserPersistence) IsValidUsernameForUser(ctx context.Context, s string, u uuid.UUID) (bool, error) {
	return true, nil
}

func (f fakeUserPersistence) GetUsers(ctx context.Context, request *shared.PaginationRequest) ([]*models.User, error) {
	return []*models.User{f.mock}, nil
}

func (f fakeUserPersistence) CreateUser(ctx context.Context, user *models.User) error {
	return nil
}

func (f fakeUserPersistence) UpdateUser(ctx context.Context, user *models.User) error {
	return nil
}

func (f fakeUserPersistence) DeleteUser(ctx context.Context, u uuid.UUID) error {
	return nil
}

func TestSimpleJwtTest(t *testing.T) {

	user := &models.User{
		ID:            uuid.New(),
		Administrator: true,
		Name:          "John",
	}

	handler := JwtHandler{
		UserPersistence:   fakeUserPersistence{user},
		TimeProvider:      providers.UtcTimeProvider{},
		RefreshExpiration: 20 * time.Second,
		Expiration:        2100 * time.Second,
		Key:               []byte("test"),
		RefreshKey:        []byte("tset"),
	}

	jwt, err := handler.CreateJwt(user)
	if err != nil {
		t.Error(err)
	}

	validateJwt, err := handler.ValidateJwt(jwt.Token)
	if err != nil {
		t.Error(err)
	}

	if validateJwt.Id != user.ID || validateJwt.Username != user.Name || validateJwt.Admin != user.Administrator {
		t.Fail()
	}

	_, err = handler.RefreshToken(context.Background(), jwt.RefreshToken)
	if err != nil {
		t.Error(err)
	}
}
