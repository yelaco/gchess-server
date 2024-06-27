package api

import (
	"encoding/json"
	"net/http"

	"github.com/yelaco/go-chess-server/internal/database"
)

func handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
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

	newUser, err := database.CreateUser(params.Username, params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, userResponse{
		PlayerID: newUser.PlayerID,
		Username: newUser.Username,
	})
}
