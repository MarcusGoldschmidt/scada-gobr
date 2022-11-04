package pkg

import (
	"github.com/MarcusGoldschmidt/scadagobr/pkg/models"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/server"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

func GetViewsHandler(s *Scadagobr, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	views, err := s.viewPersistence.GetAllViews(r.Context())
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	s.respondJsonOk(ctx, w, views)
}

func GetViewByIdHandler(s *Scadagobr, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// TODO: verify if user has permission for `the view
	viewId, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	view, err := s.viewPersistence.GetViewById(r.Context(), viewId)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	s.respondJsonOk(ctx, w, view)
}

type ViewRequest struct {
	Name           string                  `json:"name"`
	ViewComponents []*ViewComponentRequest `json:"viewComponents"`
}

type ViewComponentRequest struct {
	Id     uuid.UUID       `json:"id"`
	X      int             `json:"x"`
	Y      int             `json:"y"`
	ViewId uuid.UUID       `json:"viewId"`
	Type   models.ViewType `json:"type"`
	Data   map[string]any  `json:"data"`
}

func CreateViewHandler(s *Scadagobr, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := server.ValidateFromBody[ViewRequest](r)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	view := &models.View{
		Id:   uuid.New(),
		Name: body.Name,
	}

	err = s.viewPersistence.CreateView(r.Context(), view)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	viewComponents := make([]*models.ViewComponent, len(body.ViewComponents))

	for i, component := range body.ViewComponents {
		viewComponents[i] = &models.ViewComponent{
			Id:       component.Id,
			ViewType: component.Type,
			ViewId:   view.Id,
			Position: *models.NewPosition(component.X, component.Y),
			Data:     component.Data,
		}
	}

	err = s.viewPersistence.AttachViewComponents(r.Context(), viewComponents...)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func UpdateViewHandler(s *Scadagobr, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	viewId, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	body, err := server.ValidateFromBody[ViewRequest](r)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}
	view := &models.View{
		Id:   viewId,
		Name: body.Name,
	}

	err = s.viewPersistence.UpdateView(r.Context(), view)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	viewComponents := make([]*models.ViewComponent, len(body.ViewComponents))

	for i, component := range body.ViewComponents {
		viewComponents[i] = &models.ViewComponent{
			Id:       component.Id,
			ViewType: component.Type,
			ViewId:   viewId,
			Position: *models.NewPosition(component.X, component.Y),
			Data:     component.Data,
		}
	}

	err = s.viewPersistence.AttachViewComponents(r.Context(), viewComponents...)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteViewHandler(s *Scadagobr, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	viewId, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	err = s.viewPersistence.DeleteView(r.Context(), viewId)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteViewComponentHandler(s *Scadagobr, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	viewId, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	componentId, err := uuid.Parse(mux.Vars(r)["componentId"])
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	err = s.viewPersistence.DeleteViewComponent(r.Context(), viewId, componentId)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
