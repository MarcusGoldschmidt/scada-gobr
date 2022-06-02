package pkg

import (
	"encoding/json"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/auth"
	"io/ioutil"
	"net/http"
)

type loginRequest struct {
	Name     string
	Password string
}

func LoginHandler(s *Scadagobr, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request loginRequest
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.respondError(w, err)
		return
	}

	err = json.Unmarshal(reqBody, &request)
	if err != nil {
		s.respondError(w, err)
		return
	}

	user, err := s.userPersistence.GetUserByName(ctx, request.Name)
	if err != nil {
		s.respondError(w, err)
		return
	}

	success, err := auth.ValidatePassword(request.Password, user.PasswordHash)
	if err != nil {
		s.respondError(w, err)
		return
	}

	if !success {
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	response, err := s.JwtHandler.CreateJwt(user)
	if err != nil {
		s.respondError(w, err)
		return
	}

	s.respondJsonOk(w, response)
}

func RefreshTokenHandler(s *Scadagobr, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request map[string]string
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.respondError(w, err)
		return
	}

	err = json.Unmarshal(reqBody, &request)
	if err != nil {
		s.respondError(w, err)
		return
	}

	response, err := s.JwtHandler.RefreshToken(ctx, request["refreshToken"])
	if err != nil {
		s.respondError(w, err)
		return
	}

	s.respondJsonOk(w, response)
}

func WhoAmIHandler(s *Scadagobr, w http.ResponseWriter, r *http.Request) {

	claims, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		s.respondError(w, err)
		return
	}

	s.respondJsonOk(w, claims)
}
