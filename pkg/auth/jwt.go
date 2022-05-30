package auth

import (
	"errors"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/models"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/persistence"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/providers"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"time"
)

var ErrUnauthorized error = errors.New("Unauthorized")
var ContexClaimsKey string = "CLAIMS"

type Claims struct {
	Id       uuid.UUID
	Username string
	Admin    bool
	jwt.StandardClaims
}

type JwtHandler struct {
	Expiration        time.Duration
	RefreshExpiration time.Duration
	RefreshKey        []byte
	Key               []byte

	TimeProvider    providers.TimeProvider
	UserPersistence persistence.UserPersistence
}

func (handler *JwtHandler) ValidateJwt(token string) (*Claims, error) {

	claims := &Claims{}

	tokenParse, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return handler.Key, nil
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
		return handler.RefreshKey, nil
	})
	if err != nil {
		return nil, nil, err
	}

	if !token.Valid {
		return nil, nil, ErrUnauthorized
	}

	user, err := handler.UserPersistence.GetUser(claims.Id)
	if err != nil {
		return nil, nil, err
	}

	return handler.CreateJwt(user)
}

// CreateJwt create normal token que refresh token
func (handler *JwtHandler) CreateJwt(user *models.User) (*string, *string, error) {
	expiration := handler.TimeProvider.GetCurrentTime().Add(handler.Expiration)

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

	tokenString, err := token.SignedString(handler.Key)

	if err != nil {
		return nil, nil, err
	}

	expiration = handler.TimeProvider.GetCurrentTime().Add(handler.RefreshExpiration)
	refreshClaims := jwt.StandardClaims{
		ExpiresAt: expiration.Unix(),
		Id:        uuid.NewString(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	refreshTokenString, err := refreshToken.SignedString(handler.RefreshKey)

	return &tokenString, &refreshTokenString, nil
}
