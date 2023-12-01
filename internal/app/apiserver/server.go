package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/artemiyKew/http-rest-api/internal/app/model"
	"github.com/artemiyKew/http-rest-api/internal/app/store"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
)

const (
	ctxKeyUser ctxKey = iota
	ctxKeyRequestID
)

var (
	jwtKey            = []byte("secret-key")
	errNotFoundHeader = errors.New("header not found")
	errInvalidToken   = errors.New("invalid token")
	errTokenExpired   = errors.New("token expired")
)

type ctxKey int8

type server struct {
	router       *mux.Router
	logger       *zap.Logger
	store        store.Store
	sessionStore sessions.Store
}

func newServer(store store.Store, sessionStore sessions.Store) *server {
	s := &server{
		router:       mux.NewRouter(),
		logger:       zap.NewNop(),
		store:        store,
		sessionStore: sessionStore,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))

	s.router.HandleFunc("/sign-up", s.handleUsersCreate()).Methods("POST")
	s.router.HandleFunc("/sign-in", s.handleSessionsCreate()).Methods("POST")

	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authUser)
	private.HandleFunc("/whoami", s.handleWhoami()).Methods("GET")
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger
		msg := fmt.Sprintf("remote_addr=%s request_id=%s", r.RemoteAddr, r.Context().
			Value(ctxKeyRequestID))
		logger.Info(fmt.Sprintf("started %s %s \t %s", r.Method, r.RequestURI, msg))

		start := time.Now()

		rw := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(rw, r)
		logger.Info(fmt.Sprintf("completed with %d % s in %v \t %s",
			rw.code,
			http.StatusText(rw.code),
			time.Since(start),
			msg),
		)
	})
}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *server) authUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := s.validateToken(w, r)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			s.error(w, r, http.StatusInternalServerError, errors.New("cannot parse data"))
			return
		}
		exp := claims["exp"].(float64)
		if int64(exp) < time.Now().Local().Unix() {
			s.error(w, r, http.StatusInternalServerError, errTokenExpired)
			return
		}
		u, err := s.store.User().FindByID(int(claims["sub"].(float64)))
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))
	})
}

func (s *server) handleUsersCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}
		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()
		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) handleSessionsCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByEmail(req.Email)
		if err != nil || !u.ComparePasswords(req.Password) {
			s.error(w, r, http.StatusUnauthorized, store.ErrEmailORPasswordInvalid)
			return
		}

		token, err := generateJWT(u)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		w.Header().Set("Token", token)
		s.respond(w, r, http.StatusOK, token)
	}
}

func (s *server) validateToken(w http.ResponseWriter, r *http.Request) (*jwt.Token, error) {
	if r.Header["Token"] == nil {
		s.error(w, r, http.StatusBadRequest, errNotFoundHeader)
		return nil, errNotFoundHeader
	}

	token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("cannot parse")
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if token == nil {
		return nil, errInvalidToken
	}

	return token, nil
}

func (s *server) handleWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, r.Context().Value(ctxKeyUser).(*model.User))
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			return
		}
	}
}
