package http

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jawr/catcher/service/internal/catcher"
)

const (
	defaultTimeout time.Duration = 100
	defaultAddr    string        = "localhost:8000"
)

type Config struct {
	Address      string        `yaml:"address"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	SiteRoot     string        `yaml:"site_root"`
}

// Server wraps a store and an http.Server to resolve API requests
type Server struct {
	store    catcher.Store
	httpd    *http.Server
	siteRoot string
}

// NewServer validates and configures a Server
func NewServer(config Config, store catcher.Store) (*Server, error) {
	if _, err := os.Stat(filepath.Join(config.SiteRoot, "index.html")); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("%w: %s", catcher.ErrInvalid, err)
	}

	s := Server{
		store:    store,
		siteRoot: config.SiteRoot,
	}

	router := mux.NewRouter()

	api := router.PathPrefix("/api/v1/").Subrouter()

	api.HandleFunc("/subscribe", s.handleSubscribe())
	api.HandleFunc("/random", s.handleRandomEmail).Methods("GET")

	router.PathPrefix("/").HandlerFunc(s.handleSPA)

	if len(config.Address) == 0 {
		config.Address = defaultAddr
	}

	if config.ReadTimeout == 0 {
		config.ReadTimeout = defaultTimeout
	}

	if config.ReadTimeout == 0 {
		config.ReadTimeout = defaultTimeout
	}

	s.httpd = &http.Server{
		Handler:      handlers.LoggingHandler(os.Stdout, router),
		Addr:         config.Address,
		WriteTimeout: config.WriteTimeout * time.Second,
		ReadTimeout:  config.ReadTimeout * time.Second,
	}

	return &s, nil
}

// ListenAndServe starts the underlying server
func (s *Server) ListenAndServe() error {
	return s.httpd.ListenAndServe()
}

// Close closes the underlying server
func (s *Server) Close() error {
	return s.httpd.Close()
}
