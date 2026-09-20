package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/gmclean3107/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	queries        *database.Queries
}

func main() {
	const fileRoot = "."
	const port = "8080"

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	handler := http.NewServeMux()
	apiConfig := apiConfig{
		queries: dbQueries,
	}

	handler.Handle("/app/", apiConfig.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(fileRoot)))))

	handler.HandleFunc("GET /api/healthz", handlerReadiness)
	handler.HandleFunc("POST /api/validate_chirp", handlerValidatePost)
	handler.HandleFunc("POST /api/users", apiConfig.handlerCreateUser)

	handler.HandleFunc("GET /admin/metrics", apiConfig.handlerMetrics)
	handler.HandleFunc("POST /admin/reset", apiConfig.handlerResetApi)

	server := http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	log.Printf("Serving files from %s on port: %s\n", fileRoot, port)

	log.Fatal(server.ListenAndServe())

}
