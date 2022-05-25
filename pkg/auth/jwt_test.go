package auth

import (
	"github.com/google/uuid"
	"scadagobr/pkg/models"
	"scadagobr/pkg/providers"
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
		userPersistence:   fakeUserPersistence{user},
		timeProvider:      providers.UtcTimeProvider{},
		refreshExpiration: 20 * time.Second,
		expiration:        2100 * time.Second,
		key:               []byte("test"),
		refreshKey:        []byte("tset"),
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
