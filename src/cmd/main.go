package main

import (
	"flag"
	"fmt"
	"gocloud/src/router"
	"gocloud/src/server"

	"log/slog"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

var version = "v1.0.0"

func main() {
	logHandler := log.NewWithOptions(os.Stdout, log.Options{
		ReportTimestamp: true,
	})
	logger := slog.New(logHandler)

	var cfg struct {
		server  server.Config
		version string
	}

	flag.IntVar(&cfg.server.Port, "port", 3000, "API server port")
	flag.StringVar(&cfg.server.Host, "host", "localhost", "API server host")
	flag.StringVar(&cfg.server.Env, "env", "development", "Environment (development|staging|production)")
	flag.DurationVar(&cfg.server.ReadTimeout, "server-read-timeout", 5*time.Second, "Maximum duration for reading the entire request, including the body.")
	flag.DurationVar(&cfg.server.WriteTimeout, "server-write-timeout", 10*time.Second, "Maximum duration for writing the response, including the body.")
	flag.DurationVar(&cfg.server.IdleTimeout, "server-idle-timeout", 120*time.Second, "Maximum amount of time to wait for the next request when keep-alives are enabled.")
	flag.DurationVar(&cfg.server.ShutdownTimeout, "server-shutdown-timeout", 30*time.Second, "Maximum duration to wait for active connections to close during server shutdown.")

	var showAPIVersion bool
	flag.BoolVar(&showAPIVersion, "version", false, "Show API version")

	flag.Parse()
	if showAPIVersion {
		fmt.Printf("API Version: %s\n", version)
		os.Exit(0)
	}

	router := router.NewRouter(router.Config{
		Logger:          logger,
		API_Environment: cfg.server.Env,
		API_Version:     version,
	})

	srv := server.NewServer(logger, router, cfg.server)
	if err := srv.Start(); err != nil {
		logger.Error("server failed to gracefully shutdown", "error", err.Error())
		os.Exit(1)
	}
}
