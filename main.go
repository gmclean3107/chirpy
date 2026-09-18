package main

import (
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
	handler.HandleFunc("GET /api/healthz", handlerReadiness)
	handler.HandleFunc("GET /api/metrics", apiConfig.handlerMetrics)
	handler.HandleFunc("POST /api/reset", apiConfig.handlerResetMetrics)

	server := http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	log.Printf("Serving files from %s on port: %s\n", fileRoot, port)

	log.Fatal(server.ListenAndServe())

}
