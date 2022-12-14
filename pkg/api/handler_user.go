package api

import (
	"errors"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/auth"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/models"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/server"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type UserResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Email         *string   `json:"email"`
	HomeUrl       string    `json:"homeUrl"`
	Administrator bool      `json:"administrator"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func GetUsersHandler(s *ScadaApi, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	request, err := shared.NewPaginationRequest(r)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	users, err := s.UserPersistence.GetUsers(ctx, request)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	s.respondJsonOk(ctx, w, users)
}

func GetUserHandler(s *ScadaApi, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		s.respondError(ctx, w, errors.New("id must be a uuid4"))
		return
	}

	claims, err := auth.GetUserFromContext(ctx)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	if !claims.Admin {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := s.UserPersistence.GetUserById(ctx, id)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}
	s.respondJsonOk(ctx, w, user)
}

type CreateUserRequest struct {
	Name          string `validate:"required,alphaunicode"`
	Email         *string
	HomeUrl       string
	Administrator bool
	Password      string `validate:"required,gte=6"`
}

func (r *CreateUserRequest) ToUser() *models.User {
	hash, _ := auth.MakeHash(r.Password)

	return &models.User{
		ID:            uuid.New(),
		Name:          r.Name,
		Email:         r.Email,
		Administrator: r.Administrator,
		HomeUrl:       r.HomeUrl,
		PasswordHash:  hash,
	}
}

func CreateUserHandler(s *ScadaApi, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	request, err := server.ValidateFromBody[CreateUserRequest](r)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	if *request.Email != "" {
		err := server.Validate.Var(request.Email, "email")
		if err != nil {
			s.respondError(ctx, w, err)
			return
		}
	}

	userName, err := s.UserPersistence.GetUserByName(ctx, request.Name)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	if userName != nil {
		s.respondError(ctx, w, errors.New("the username already in use"))
		return
	}

	user := request.ToUser()

	err = s.UserPersistence.CreateUser(ctx, user)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	s.respondJson(ctx, w, http.StatusCreated, request)
}

type UpdateUserRequest struct {
	Name          string `validate:"required,alphaunicode"`
	Email         *string
	HomeUrl       string
	Administrator bool
	Password      string `validate:"omitempty,gte=6"`
}

func UpdateUserHandler(s *ScadaApi, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		s.respondError(ctx, w, errors.New("id must be a uuid4"))
		return
	}

	claims, err := auth.GetUserFromContext(ctx)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	request, err := server.ValidateFromBody[UpdateUserRequest](r)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	if *request.Email != "" {
		err := server.Validate.Var(request.Email, "email")
		if err != nil {
			s.respondError(ctx, w, err)
			return
		}
	}

	user, err := s.UserPersistence.GetUserById(ctx, id)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	valid, err := s.UserPersistence.IsValidUsernameForUser(ctx, request.Name, user.ID)

	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	if !valid {
		s.respondError(ctx, w, errors.New("the username already in use"))
		return
	}

	user.Name = request.Name
	user.HomeUrl = request.HomeUrl
	user.Administrator = request.Administrator
	user.Email = request.Email

	if id == claims.Id && claims.Admin {
		user.Administrator = true
	}

	if request.Password != "" {
		user.PasswordHash, _ = auth.MakeHash(request.Password)
	}

	err = s.UserPersistence.UpdateUser(ctx, user)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	s.respondJsonOk(ctx, w, map[string]any{})
}

func DeleteUserHandler(s *ScadaApi, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		s.respondError(ctx, w, errors.New("id must be a uuid4"))
		return
	}

	claims, err := auth.GetUserFromContext(ctx)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	if !claims.Admin {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if id == claims.Id {
		s.respondError(ctx, w, errors.New("cannot delete your own user"))
		return
	}

	err = s.UserPersistence.DeleteUser(ctx, id)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	s.respondJsonOk(ctx, w, map[string]any{})
}
