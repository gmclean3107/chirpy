package main

import "net/http"

func (cfg *apiConfig) handlerResetApi(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Store(0)
	err := cfg.queries.DeleteAllUsers(r.Context())

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error resetting users table", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0\n"))
	w.Write([]byte("Users table reset\n"))
}
