package api

import (
	"net/http"

	"github.com/yelaco/go-chess-server/internal/database"
)

func handlerSessionGet(w http.ResponseWriter, r *http.Request) {
	player_id := r.URL.Query().Get("player_id")
	if player_id == "" {
		respondWithError(w, http.StatusBadRequest, "Player id not included")
		return
	}

	sessionIDs, err := database.GetSessionsByPlayerID(player_id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid player id")
	}

	respondWithJSON(w, http.StatusOK, sessionIDs)
}

func handlerSessionGetFromID(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("sessionid")

	session, err := database.GetSessionByID(sessionID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid session id")
	}

	respondWithJSON(w, http.StatusOK, session)
}
