package main

import (
	"encoding/json"
	"net/http"
	"github.com/google/uuid"
	"github.com/timbdn01/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	/* accepts a webhook from polka, which will have the following JSON body:
{"event": "user.upgraded",
"data": {
"user_id": "uuid of the user that was upgraded"
}}
If the event is "user.upgraded", then we will upgrade the user to chirpy red by setting the is_chirpy_red field in the database to true. We will return a 204 status code if the upgrade was successful, or a 204 status code if the event was not recognized, and a 500 status code if there was an error upgrading the user. */
	type parameters struct {
	Event string `json:"event"`
	Data struct {
		UserID string `json:"user_id"`
	} `json:"data"`
	}
	//check the APIKey in the header matches the one in the config
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}
	if apiKey != cfg.polka_key {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	if params.Event != "user.upgraded" {
		respondWithError(w, http.StatusNoContent, "Unrecognized event", nil)
		return
	}
	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}
	err = cfg.db.UpgradeToChirpyRed(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't upgrade user", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, map[string]string{
		"message": "User upgraded to Chirpy Red successfully",
		"user_id": userID.String(),
	})
}