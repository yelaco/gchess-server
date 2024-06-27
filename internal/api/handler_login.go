package api

import (
	"encoding/json"
	"net/http"

	"github.com/yelaco/go-chess-server/internal/auth"
	"github.com/yelaco/go-chess-server/internal/database"
)

func handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := database.GetUserByUsername(params.Username)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find user with specified email")
		return
	}

	if err := auth.CheckPasswordHash(user.Password, params.Password); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Password not matched")
		return
	}

	respondWithJSON(w, http.StatusOK, userResponse{
		PlayerID: user.PlayerID,
		Username: user.Username,
	})
}
