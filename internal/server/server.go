package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

type muxRouter interface {
	HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route
	PathPrefix(tpl string) *mux.Route
}

type Controller interface {
	RegisterRoutes(r muxRouter)
}

func apiMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func NewServer(controllers []Controller) http.Handler {
	r := mux.NewRouter()

	// Ensure routes with a trailing slash get redirected
	r.StrictSlash(true)

	for _, c := range controllers {
		c.RegisterRoutes(r)
	}

	return r
}
