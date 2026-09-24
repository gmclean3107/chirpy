package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gmclean3107/chirpy/internal/database"
	"github.com/google/uuid"
)

type parameters struct {
	UserID uuid.UUID `json:"user_id"`
	Body   string    `json:"body"`
}

type returnVals struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Body == "" {
		respondWithError(w, http.StatusBadRequest, "Chirp body is required", nil)
		return
	}

	if params.UserID == uuid.Nil {
		respondWithError(w, http.StatusBadRequest, "User ID is required", nil)
		return
	}

	if !cfg.handlerValidatePost(w, &params) {
		return
	}

	res, err := cfg.queries.CreateChirp(r.Context(), database.CreateChirpParams{
		UserID: params.UserID,
		Body:   params.Body,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, returnVals{
		ID:        res.ID,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
		Body:      res.Body,
		UserID:    res.UserID,
	})

}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.queries.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch chirps", err)
		return
	}

	chirpResponses := []returnVals{}
	for _, chirp := range chirps {
		chirpResponses = append(chirpResponses, returnVals{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirpResponses)
}

func (cfg *apiConfig) handlerValidatePost(w http.ResponseWriter, params *parameters) bool {

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return false
	}

	replaceProfanity(&params.Body)
	return true

}

func replaceProfanity(chirp *string) {
	replacer := strings.NewReplacer("kerfuffle", "****", "sharbert", "****", "fornax", "****")

	split := strings.Split(*chirp, " ")

	for i := range split {
		original := split[i]
		split[i] = replacer.Replace(strings.ToLower(split[i]))
		if split[i] != "****" {
			split[i] = original
		}
	}

	*chirp = strings.Join(split, " ")
}
