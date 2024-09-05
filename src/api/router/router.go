package router

import (
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

	r := chi.NewRouter()

	return r
}
