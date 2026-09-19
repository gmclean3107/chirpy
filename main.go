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
	handler.HandleFunc("POST /api/validate_chirp", handlerValidatePost)

	handler.HandleFunc("GET /admin/metrics", apiConfig.handlerMetrics)
	handler.HandleFunc("POST /admin/reset", apiConfig.handlerResetMetrics)

	server := http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	log.Printf("Serving files from %s on port: %s\n", fileRoot, port)

	log.Fatal(server.ListenAndServe())

}
