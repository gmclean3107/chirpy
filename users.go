package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	type returnVals struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	reqBody := r.Body

	reqBytes, err := io.ReadAll(reqBody)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error reading request body to byte slice", err)
		return
	}

	params := parameters{}
	err = json.Unmarshal(reqBytes, &params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error reading request body to byte slice", err)
		return
	}

	if params.Email == "" {
		respondWithError(w, http.StatusBadRequest, "email field not provided", errors.New("No email field in JSON"))
		return
	}

	user, err := cfg.queries.CreateUser(r.Context(), params.Email)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating user", err)
		return
	}

	res := returnVals{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJSON(w, http.StatusCreated, res)
}
