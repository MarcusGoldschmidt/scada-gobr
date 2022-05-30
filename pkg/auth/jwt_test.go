package auth

import (
	"github.com/MarcusGoldschmidt/scadagobr/pkg/models"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/providers"
	"github.com/google/uuid"
	"testing"
	"time"
)

type fakeUserPersistence struct {
	mock *models.User
}

func (f fakeUserPersistence) GetUser(uuid uuid.UUID) (*models.User, error) {
	return f.mock, nil
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

	jwt, refresh, err := handler.CreateJwt(user)
	if err != nil {
		t.Error(err)
	}

	validateJwt, err := handler.ValidateJwt(*jwt)
	if err != nil {
		t.Error(err)
	}

	if validateJwt.Id != user.ID || validateJwt.Username != user.Name || validateJwt.Admin != user.Administrator {
		t.Fail()
	}

	_, _, err = handler.RefreshToken(*refresh)
	if err != nil {
		t.Error(err)
	}
}
