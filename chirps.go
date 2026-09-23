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

	cfg.handlerValidatePost(w, r, params)

}

func (cfg *apiConfig) handlerValidatePost(w http.ResponseWriter, r *http.Request, params parameters) {

	type returnVals struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	replaceProfanity(&params.Body)

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
