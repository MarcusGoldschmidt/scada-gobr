package server

import (
	"net/http"
	"sync"
)

type Router struct {
	sync.RWMutex
	match map[string]http.Handler
}

func (r *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	match := request.URL.Path

	var handler http.Handler

	r.RLock()
	for key, value := range r.match {
		if key == match {
			handler = value
		}
	}
	r.RUnlock()

	if handler != nil {
		handler.ServeHTTP(writer, request)
		return
	}

	http.NotFound(writer, request)
}

func NewRouter() *Router {
	return &Router{sync.RWMutex{}, make(map[string]http.Handler)}
}

func (r *Router) VerifyMatch(match string) (http.Handler, bool) {
	r.RLock()
	defer r.RUnlock()

	data := r.match[match]
	return data, data != nil
}

func (r *Router) AddMatch(match string, handler http.Handler) {
	r.Lock()
	defer r.Unlock()
	r.match[match] = handler
}

func (r *Router) RemoveMatch(match string) {
	r.Lock()
	defer r.Unlock()

	delete(r.match, match)
}
