package auth

import (
	"context"
	"errors"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/models"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/persistence"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/providers"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"time"
)

var ErrUnauthorized error = errors.New("Unauthorized")
var ContextClaimsKey string = "CLAIMS"

func GetUserFromContext(ctx context.Context) (*Claims, error) {
	value := ctx.Value(ContextClaimsKey)

	claims, ok := value.(*Claims)
	if ok {
		return claims, nil
	}

	return nil, errors.New("unable to receive claims of the current user")
}

type Claims struct {
	Id       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Admin    bool      `json:"admin"`
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

type JwtResponse struct {
	Token                  string `json:"token"`
	TokenExpiration        int64  `json:"tokenExpiration"`
	RefreshToken           string `json:"refreshToken"`
	RefreshTokenExpiration int64  `json:"refreshTokenExpiration"`
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

func (handler *JwtHandler) RefreshToken(refreshToken string) (*JwtResponse, error) {

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return handler.RefreshKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrUnauthorized
	}

	user, err := handler.UserPersistence.GetUserById(claims.Id)
	if err != nil {
		return nil, err
	}

	return handler.CreateJwt(user)
}

// CreateJwt create normal token que refresh token
func (handler *JwtHandler) CreateJwt(user *models.User) (*JwtResponse, error) {
	expirationToken := handler.TimeProvider.GetCurrentTime().Add(handler.Expiration)

	claims := &Claims{
		Id:       user.ID,
		Username: user.Name,
		Admin:    user.Administrator,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationToken.Unix(),
			Id:        uuid.NewString(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(handler.Key)

	if err != nil {
		return nil, err
	}

	expirationRefreshToken := handler.TimeProvider.GetCurrentTime().Add(handler.RefreshExpiration)
	refreshClaims := jwt.StandardClaims{
		ExpiresAt: expirationRefreshToken.Unix(),
		Id:        uuid.NewString(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	refreshTokenString, err := refreshToken.SignedString(handler.RefreshKey)

	response := &JwtResponse{
		Token:                  tokenString,
		TokenExpiration:        expirationToken.Unix(),
		RefreshToken:           refreshTokenString,
		RefreshTokenExpiration: expirationRefreshToken.Unix(),
	}

	return response, nil
}
