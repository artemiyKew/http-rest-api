package apiserver

import (
	"encoding/json"
	"net/http"

	"github.com/artemiyKew/http-rest-api/internal/app/model"
	"github.com/artemiyKew/http-rest-api/internal/app/store"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
)

const (
	sessionName = "artemiyKew"
)

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
	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods("POST")
	s.router.HandleFunc("/sessions", s.handleSessionsCreate()).Methods("POST")
	s.router.HandleFunc("/user", s.handleUserData()).Methods("Get")
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

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["user_id"] = u.ID
		if err := s.sessionStore.Save(r, w, session); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleUserData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			if err == http.ErrNoCookie {
				s.error(w, r, http.StatusUnauthorized, err)
				return
			}
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		auth, ok := session.Values["user_id"].(int)
		if !ok {
			s.error(w, r, http.StatusForbidden, err)
			return
		}

		s.respond(w, r, http.StatusOK, auth)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// func (s *server) configureLogger() error {
// 	filename := "logs.log"
// 	config := zap.NewProductionEncoderConfig()
// 	config.EncodeTime = zapcore.ISO8601TimeEncoder

// 	// Create file and console encoders
// 	fileEncoder := zapcore.NewJSONEncoder(config)
// 	consoleEncoder := zapcore.NewConsoleEncoder(config)

// 	logFile, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		return err
// 	}

// 	// Create writers for file and console
// 	fileWriter := zapcore.AddSync(logFile)
// 	consoleWriter := zapcore.AddSync(os.Stdout)

// 	// Set the log level
// 	defaultLogLevel := zapcore.InfoLevel

// 	// Create cores for writing to the file and console
// 	fileCore := zapcore.NewCore(fileEncoder, fileWriter, defaultLogLevel)
// 	consoleCore := zapcore.NewCore(consoleEncoder, consoleWriter, defaultLogLevel)

// 	// Combine cores
// 	core := zapcore.NewTee(fileCore, consoleCore)
// 	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

// 	s.logger = logger
// 	return nil
// }
