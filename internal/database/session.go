package database

import (
	"encoding/json"
	"log"
)

type Session struct {
	SessionID string   `json:"session_id"`
	Player1ID string   `json:"player1_id"`
	Player2ID string   `json:"player2_id"`
	Moves     []string `json:"moves"`
}

func GetSessionByID(sessionID string) (Session, error) {
	var session Session
	query := `SELECT session_id, player1_id, player2_id, moves FROM sessions WHERE session_id = $1`
	row := db.QueryRow(query, sessionID)

	var moveJSON string
	err := row.Scan(&session.SessionID, &session.Player1ID, &session.Player2ID, &moveJSON)
	if err != nil {
		return Session{}, err
	}

	err = json.Unmarshal([]byte(moveJSON), &session.Moves)
	if err != nil {
		return Session{}, err
	}

	return session, nil
}

func GetSessionsByPlayerID(playerID string) ([]Session, error) {
	var sessions []Session

	query := `SELECT session_id, player1_id, player2_id, moves FROM sessions WHERE player1_id = $1 OR player2_id = $1`
	rows, err := db.Query(query, playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var session Session
		var movesJSON string
		err := rows.Scan(&session.SessionID, &session.Player1ID, &session.Player2ID, &movesJSON)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(movesJSON), &session.Moves); err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}

func InsertSession(sessionID, player1ID, player2ID string, moves []string) (Session, error) {
	movesJSON, err := json.Marshal(moves)
	if err != nil {
		log.Fatal(err)
	}

	ist, err := db.Prepare("INSERT INTO sessions (session_id, player1_id, player2_id, moves) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return Session{}, err
	}
	defer ist.Close()

	_, err = ist.Exec(sessionID, player1ID, player2ID, movesJSON)
	if err != nil {
		return Session{}, err
	}

	return Session{
		SessionID: sessionID,
		Player1ID: player1ID,
		Player2ID: player2ID,
		Moves:     moves,
	}, nil
}
