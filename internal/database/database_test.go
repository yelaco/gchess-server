package database

import (
	"fmt"
	"testing"
)

func TestUser(t *testing.T) {
	InitDB()
	if db == nil {
		t.Error("nil db")
		return
	}

	newUser, err := CreateUser("tester2", "password")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(newUser)

	user, err := GetUserByUsername(newUser.Username)
	if err != nil {
		t.Error(err)
		return
	}
	if user.PlayerID != newUser.PlayerID {
		t.Errorf("get user: got %v, want %v", user.PlayerID, newUser.PlayerID)
		return
	}
	fmt.Println(user)

	CloseDB()
}

func TestSession(t *testing.T) {
	InitDB()
	if db == nil {
		t.Error("nil db")
	}

	newSession, err := InsertSession("1234", "fd9a179f-c035-4e50-82f5-5d1efc844316", "0046bb25-3f06-44f8-84e2-d84e2fff42e9", []string{"e2-e4"})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(newSession)

	sessions, err := GetSessionsByPlayerID(newSession.Player1ID)
	if err != nil {
		t.Error(err)
		return
	}

	session, err := GetSessionByID(sessions[0].SessionID)
	if err != nil {
		t.Error(err)
		return
	}
	if session.Player1ID != newSession.Player1ID {
		t.Errorf("get session: got %v, want %v", session.Player1ID, session.Player2ID)
	}
	fmt.Println(session)

	CloseDB()
}
