package router

import (
	"gocloud/src/health"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Config struct {
	API_Environment string
	API_Version     string
	Logger          *slog.Logger
}

func NewRouter(cfg Config) http.Handler {
	router := chi.NewRouter()
	health.RegisterHandlers(router)
	return router
}
