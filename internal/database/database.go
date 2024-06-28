package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/yelaco/go-chess-server/pkg/logging"
	"go.uber.org/zap"
)

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("postgres", "user=server password=chessserver dbname=chess sslmode=disable")
	if err != nil {
		logging.Fatal("database connection failure", zap.Error(err))
	}

	// Ping database to verify the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	logging.Info("database connected")

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
}

func CloseDB() {
	db.Close()
}
