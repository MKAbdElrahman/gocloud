package server_test

import (
	"gocloud/src/server"
	"log/slog"
	"net"
	"net/http"
	"os"
	"testing"
	"time"
)

// mockHandler is a basic HTTP handler for testing purposes.
func mockHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// TestServerStartShutdown tests that the server can start and shut down correctly.
func TestServerStartShutdown(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	// Create a logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Define server configuration

	port, _ := getAvailablePort()
	cfg := server.Config{
		Port:            port,
		Host:            "127.0.0.1",
		Env:             "test",
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    5 * time.Second,
		IdleTimeout:     5 * time.Second,
		ShutdownTimeout: 5 * time.Second,
	}

	// Create a new server instance with a mock handler
	srv := server.NewServer(logger, http.HandlerFunc(mockHandler), cfg)

	// Run the server in a goroutine
	go func() {
		if err := srv.Start(); err != nil {
			t.Errorf("Server failed to start: %v", err)
		}
	}()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Check if the server responds to HTTP requests
	resp, err := http.Get("http://" + srv.Address())
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}

	// Manually shut down the server
	if err := srv.Shutdown(); err != nil {
		t.Errorf("Failed to shut down server: %v", err)
	}
}

func getAvailablePort() (int, error) {
	// Listen on a random port with "tcp"
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()

	// Get the listener's assigned address, which contains the available port
	addr := listener.Addr().(*net.TCPAddr)
	return addr.Port, nil
}
