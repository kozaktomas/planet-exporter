package main

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// NewRouter creates a new HTTP router
func NewRouter(solarSystem *SolarSystem, monitor *Monitor) *http.ServeMux {
	router := http.NewServeMux()
	router.Handle("/metrics", metricsHandler(solarSystem, monitor))
	router.Handle("/", homepageHandler())

	return router
}

// StartServer starts HTTP server
// It listens for SIGINT and SIGTERM signals and gracefully stops the server
func StartServer(router *http.ServeMux, port int) {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("listen: %s\n", err)
		}
	}()
	log.Printf("Server Started on port %d", port)

	<-done
	log.Printf("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	log.Printf("Server Exited Properly")
}

// metricsHandler returns HTTP handler for metrics endpoint
func metricsHandler(solarSystem *SolarSystem, monitor *Monitor) http.Handler {
	return refreshMiddleWare(promhttp.HandlerFor(
		monitor.Registry,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
			Registry:          monitor.Registry,
		},
	), solarSystem)
}

// refreshMiddleWare is a middleware that initiate recalculation of object positions
func refreshMiddleWare(next http.Handler, solarSystem *SolarSystem) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := solarSystem.recalculatePositions(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("could not load data"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// homepageHandler returns HTTP handler for homepage
func homepageHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
			<head><title>Planet exporter</title></head>
			<body>
			<h1>Planet exporter</h1>
			<p><a href="/metrics">Metrics</a></p>
			</body>
			</html>`))
	})
}
