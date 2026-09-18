package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func main() {
	const fileRoot = "."
	const port = "8080"

	handler := http.NewServeMux()
	apiConfig := apiConfig{}

	handler.Handle("/app/", apiConfig.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(fileRoot)))))
	handler.HandleFunc("/healthz", handlerReadiness)
	handler.HandleFunc("/metrics", apiConfig.handlerMetrics)
	handler.HandleFunc("/reset", apiConfig.handlerResetMetrics)

	server := http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	err := server.ListenAndServe()

	if err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(fmt.Appendf(nil, "Hits: %v", cfg.fileServerHits.Load()))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerResetMetrics(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Store(0)
}
