package health

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterHandlers(mux chi.Router) {
	mux.Get("/health", func(w http.ResponseWriter, r *http.Request) {
	})
}
