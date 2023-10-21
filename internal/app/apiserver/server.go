package apiserver

import (
	"net/http"

	"github.com/artemiyKew/http-rest-api/internal/app/store"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type server struct {
	router *mux.Router
	logger *zap.Logger
	store  *store.Store
}

func newServer(store store.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: zap.NewNop(),
		store:  &store,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods("POST")
}

func (s *server) handleUsersCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
