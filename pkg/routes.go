package pkg

import (
	"net/http"
)

type RequestHandlerFunction func(scada *Scadagobr, w http.ResponseWriter, r *http.Request)

func (s *Scadagobr) handleRequest(function RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		function(s, w, r)
	}
}

func (s *Scadagobr) get(path string, f RequestHandlerFunction) {
	s.router.HandleFunc(path, s.handleRequest(f)).Methods("GET")
}

// Post wraps the router for POST method
func (s *Scadagobr) post(path string, f RequestHandlerFunction) {
	s.router.HandleFunc(path, s.handleRequest(f)).Methods("POST")
}

// Put wraps the router for PUT method
func (s *Scadagobr) put(path string, f RequestHandlerFunction) {
	s.router.HandleFunc(path, s.handleRequest(f)).Methods("PUT")
}

// Delete wraps the router for DELETE method
func (s *Scadagobr) delete(path string, f RequestHandlerFunction) {
	s.router.HandleFunc(path, s.handleRequest(f)).Methods("DELETE")
}

func (s *Scadagobr) setRouters() {
	s.setupLogs()
	s.setupCors()
	s.setupProviders()

	s.get("/api/healthcheck", HealthCheckHandler)

	// Auth
	s.post("/api/v1/auth/login", LoginHandler)
	s.post("/api/v1/auth/refresh-token", RefreshTokenHandler)
	s.get("/api/v1/auth/who-am-i", s.jwtMiddleware(WhoAmIHandler))

	// Users
	s.get("/api/v1/user", s.authAndIsAdminMiddleware(GetUsersHandler))
	s.get("/api/v1/user/{id}", s.authAndIsAdminMiddleware(GetUserHandler))
	s.post("/api/v1/user", s.authAndIsAdminMiddleware(CreateUserHandler))
	s.put("/api/v1/user/{id}", s.authAndIsAdminMiddleware(UpdateUserHandler))
	s.delete("/api/v1/user/{id}", s.authAndIsAdminMiddleware(DeleteUserHandler))

	// Sql
	s.get("/api/v1/sql/drivers", s.authAndIsAdminMiddleware(GetDriversHandler))

	// DataSources
	s.get("/api/v1/datasource/types", s.authAndIsAdminMiddleware(GetDataSourceTypesHandler))
	s.get("/api/v1/datasource/runtime", s.authAndIsAdminMiddleware(GetDataSourcesRuntime))
	s.get("/api/v1/datasource", s.authAndIsAdminMiddleware(GetDataSourcesHandler))
	s.get("/api/v1/datasource/{id}", s.authAndIsAdminMiddleware(GetDataSourceByIdHandler))
	s.post("/api/v1/datasource", s.authAndIsAdminMiddleware(CreateDataSourceHandler))
	s.put("/api/v1/datasource/{id}", s.authAndIsAdminMiddleware(EditDataSourceHandler))
	s.delete("/api/v1/datasource/{id}", s.authAndIsAdminMiddleware(DeleteDataSourceHandler))

	// DataPoints
	s.get("/api/v1/datapoint/{id}", s.authAndIsAdminMiddleware(GetDataPointByIdHandler))
	s.get("/api/v1/datasource/{id}/datapoint", s.authAndIsAdminMiddleware(GetDataPointsHandler))
	s.post("/api/v1/datasource/{id}/datapoint", s.authAndIsAdminMiddleware(CreateDataPointHandler))
	s.put("/api/v1/datasource/{id}/datapoint/{dataPointId}", s.authAndIsAdminMiddleware(EditDataPointHandler))
	s.delete("/api/v1/datasource/{id}/datapoint/{dataPointId}", s.authAndIsAdminMiddleware(DeleteDataPointHandler))

	// RuntimeManager
	s.get("/api/v1/runtime-manager/status", s.authAndIsAdminMiddleware(GetRuntimeMangerStatusHandler))

	// DataSeries
	s.get("/api/v1/data-series/group", s.authAndIsAdminMiddleware(GetDataSeriesByGroup))

	// Websocket
	// TODO: Add authentication for websocket
	s.get("/ws/v1/view/component/time-series/{id}", GetWsTimeSeriesViewComponent)

	// Views
	s.get("/api/v1/view", s.authAndIsAdminMiddleware(GetViewsHandler))
	s.get("/api/v1/view/{id}", s.authAndIsAdminMiddleware(GetViewByIdHandler))
	s.post("/api/v1/view", s.authAndIsAdminMiddleware(CreateViewHandler))
	s.put("/api/v1/view/{id}", s.authAndIsAdminMiddleware(UpdateViewHandler))
	s.delete("/api/v1/view/{id}", s.authAndIsAdminMiddleware(DeleteViewHandler))
	s.delete("/api/v1/view/{id}/component/{componentId}", s.authAndIsAdminMiddleware(DeleteViewComponentHandler))
}
