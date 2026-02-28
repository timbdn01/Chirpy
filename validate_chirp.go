package main

import (
	"encoding/json"
	"net/http"
	"strings"
)
func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Cleaned_Body string `json:"cleaned_body"`
		Valid        bool   `json:"valid"`
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

	respondWithJSON(w, http.StatusOK, returnVals{
		Cleaned_Body: cleaned_body,
		Valid:        true,
	})

	
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