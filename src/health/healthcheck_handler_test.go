package health_test

import (
	"gocloud/src/health"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/matryer/is"
)

func TestHealthCheckHandler(t *testing.T) {
	t.Run("returns status code 200", func(t *testing.T) {
		is := is.New(t)
		router := chi.NewMux()
		health.RegisterHandlers(router)
		code, _, _ := makeGetRequest(router, "/health")
		is.Equal(code, http.StatusOK)
	})
}

func makeGetRequest(handler http.Handler, target string) (int, http.Header, string) {
	request := httptest.NewRequest(http.MethodGet, target, nil)
	responseRecorder := httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, request)
	response := responseRecorder.Result()
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	return response.StatusCode, response.Header, string(bodyBytes)
}
