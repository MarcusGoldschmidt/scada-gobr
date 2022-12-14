package api

import (
	"encoding/json"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/auth"
	"io/ioutil"
	"net/http"
)

type loginRequest struct {
	Username string
	Password string
}

func LoginHandler(s *ScadaApi, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request loginRequest
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	err = json.Unmarshal(reqBody, &request)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	user, err := s.UserPersistence.GetUserByName(ctx, request.Username)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	success := auth.ValidatePassword(request.Password, user.PasswordHash)

	if !success {
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	response, err := s.JwtHandler.CreateJwt(user)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	s.respondJsonOk(ctx, w, response)
}

func RefreshTokenHandler(s *ScadaApi, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request map[string]string
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	err = json.Unmarshal(reqBody, &request)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	response, err := s.JwtHandler.RefreshToken(ctx, request["refreshToken"])
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	s.respondJsonOk(ctx, w, response)
}

func WhoAmIHandler(s *ScadaApi, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	s.respondJsonOk(ctx, w, claims)
}
