package auth

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"scadagobr/pkg/models"
	"scadagobr/pkg/persistence"
	"scadagobr/pkg/providers"
	"time"
)

var ErrUnauthorized error = errors.New("Unauthorized")

type Claims struct {
	Id       uuid.UUID
	Username string
	Admin    bool
	jwt.StandardClaims
}

type JwtHandler struct {
	expiration        time.Duration
	refreshExpiration time.Duration
	refreshKey        []byte
	key               []byte

	timeProvider    providers.TimeProvider
	userPersistence persistence.UserPersistence
}

func (handler *JwtHandler) ValidateJwt(token string) (*Claims, error) {

	claims := &Claims{}

	tokenParse, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return handler.key, nil
	})
	if err != nil {
		return nil, err
	}

	if !tokenParse.Valid {
		return nil, ErrUnauthorized
	}

	return claims, nil
}

func (handler *JwtHandler) RefreshToken(refreshToken string) (*string, *string, error) {

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return handler.refreshKey, nil
	})
	if err != nil {
		return nil, nil, err
	}

	if !token.Valid {
		return nil, nil, ErrUnauthorized
	}

	user, err := handler.userPersistence.GetUser(claims.Id)
	if err != nil {
		return nil, nil, err
	}

	return handler.CreateJwt(user)
}

// CreateJwt create normal token que refresh token
func (handler *JwtHandler) CreateJwt(user *models.User) (*string, *string, error) {
	expiration := handler.timeProvider.GetCurrentTime().Add(handler.expiration)

	claims := &Claims{
		Id:       user.ID,
		Username: user.Name,
		Admin:    user.Administrator,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration.Unix(),
			Id:        uuid.NewString(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(handler.key)

	if err != nil {
		return nil, nil, err
	}

	expiration = handler.timeProvider.GetCurrentTime().Add(handler.refreshExpiration)
	refreshClaims := jwt.StandardClaims{
		ExpiresAt: expiration.Unix(),
		Id:        uuid.NewString(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	refreshTokenString, err := refreshToken.SignedString(handler.refreshKey)

	return &tokenString, &refreshTokenString, nil
}
