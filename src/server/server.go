package server

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type Config struct {
	Port            int
	Host            string
	Env             string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
	config     Config
}

// NewServer creates a new Server instance.
func NewServer(logger *slog.Logger, router http.Handler, cfg Config) *Server {
	srv := &http.Server{
		Addr:         net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)),
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		WriteTimeout: cfg.WriteTimeout,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	return &Server{
		httpServer: srv,
		logger:     logger,
		config:     cfg,
	}
}
func (s *Server) Address() string {
	return s.httpServer.Addr
}

// Start starts the HTTP server and listens for incoming requests.
func (s *Server) Start() error {
	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit
		s.logger.Info("received signal", "signal", sig.String())
		s.logger.Info("shutting down server", "addr", s.httpServer.Addr)

		shutdownError <- s.Shutdown()
	}()

	s.logger.Info("starting server", "addr", s.httpServer.Addr, "env", s.config.Env)
	err := s.httpServer.ListenAndServe()

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	err = <-shutdownError
	if err != nil {
		s.logger.Error("error during shutdown", "error", err)
		return err
	}

	s.logger.Info("stopped server", "addr", s.httpServer.Addr)
	return nil
}

// Shutdown gracefully shuts down the server with a given timeout.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
