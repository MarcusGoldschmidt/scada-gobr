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

	s.get("/api/healthcheck", HealthCheckHandler)

	// Auth
	s.post("/api/v1/auth/login", LoginHandler)
	s.post("/api/v1/auth/refresh-token", RefreshTokenHandler)
	s.get("/api/v1/auth/who-am-i", s.jwtMiddleware(WhoAmIHandler))

	// Users
	s.get("/api/v1/user", s.jwtMiddleware(GetUsersHandler))
	s.post("/api/v1/user", s.jwtMiddleware(CreateUserHandler))
	s.put("/api/v1/user/{id}", s.jwtMiddleware(UpdateUserHandler))
	s.delete("/api/v1/user/{id}", s.jwtMiddleware(DeleteUserHandler))
}
