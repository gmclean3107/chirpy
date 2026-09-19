package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidatePost(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	replaceProfanity(&params.Body)

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: params.Body,
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
