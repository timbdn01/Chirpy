package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/timbdn01/Chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string `json:"body"`
	UserID     uuid.UUID `json:"user_id"`
}


func profanityFilter(body string) string {
	//takes string and replaces any profane words with "****" while keeping case of characters consistent
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	body_words := strings.Split(body, " ")
	temp_body := strings.ToLower(body)
	for _, word := range profaneWords {
		temp_body = strings.ReplaceAll(temp_body, word, "****")
	}
	
	cleaned_body := []string{}
	body_parts := strings.Split(temp_body, " ")
	for i, part := range body_parts {
		if part == "****" {
			cleaned_body = append(cleaned_body, part)
		} else {
			cleaned_body = append(cleaned_body, body_words[i])
		}
	}
	return strings.Join(cleaned_body, " ")
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	//take a JSON body with "body" and "user_id" fields, validate the chirp body, and if valid, create a new chirp in the database and return status code 201 with the full chirp resource in the body. If the chirp is invalid, return status code 400 with an error message.
	type parameters struct {
		Body string `json:"body"`
		UserID string `json:"user_id"`
	}
	type returnVals struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
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
	cleaned_body := profanityFilter(params.Body)
	userID, err := uuid.Parse(params.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: cleaned_body,
		UserID: userID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, returnVals{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	})
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	//return a list of all chirps in the database, ordered by creation date, with status code 200. If there are no chirps, return an empty list.
	type returnVals struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
	}
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
		return
	}
	returnChirps := []returnVals{}
	for _, chirp := range chirps {
		returnChirps = append(returnChirps, returnVals{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID: chirp.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, returnChirps)
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	//get the chirp ID from the URL path, return the chirp with that ID if it exists, or return a 404 if it doesn't.
	type returnVals struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}
	chirpIDStr := strings.TrimPrefix(r.URL.Path, "/api/chirps/")
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}
	chirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}
	respondWithJSON(w, http.StatusOK, returnVals{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}