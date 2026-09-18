package main

import (
	"log"
	"net/http"
)

func main() {
	const fileRoot = "."
	const port = "8080"

	handler := http.NewServeMux()

	handler.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(fileRoot))))

	handler.HandleFunc("/healthz", handlerReadiness)

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
