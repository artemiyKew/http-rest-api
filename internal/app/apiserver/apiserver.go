package apiserver

import (
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ApiServer struct {
	config *Config
	logger *zap.Logger
	router *mux.Router
}

func New(config *Config) *ApiServer {
	return &ApiServer{
		config: config,
		logger: zap.NewNop(),
		router: mux.NewRouter(),
	}
}

func (s *ApiServer) configureLogger() error {
	filename := "logs.log"
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	// Create file and console encoders
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	logFile, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// Create writers for file and console
	fileWriter := zapcore.AddSync(logFile)
	consoleWriter := zapcore.AddSync(os.Stdout)

	// Set the log level
	defaultLogLevel := zapcore.DebugLevel

	// Create cores for writing to the file and console
	fileCore := zapcore.NewCore(fileEncoder, fileWriter, defaultLogLevel)
	consoleCore := zapcore.NewCore(consoleEncoder, consoleWriter, defaultLogLevel)

	// Combine cores
	core := zapcore.NewTee(fileCore, consoleCore)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	s.logger = logger
	return nil
}

func (s *ApiServer) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	s.configureRouter()

	s.logger.Info("starting api server")
	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *ApiServer) configureRouter() {
	s.router.HandleFunc("/hello", s.handleHello())
}

func (s *ApiServer) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello")
	}
}
