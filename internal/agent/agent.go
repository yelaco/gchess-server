package agent

import (
	corenet "github.com/yelaco/robinhood-chess/internal/core"
	"github.com/yelaco/robinhood-chess/internal/database"
	"github.com/yelaco/robinhood-chess/internal/logging"
)

type Agent struct {
	DBConn *database.DBConnection
}

func NewAgent() *Agent {
	dbConn := database.Connect()
	return &Agent{
		DBConn: dbConn,
	}
}

func (a *Agent) StartSocketServer() {
	go corenet.InitWebSocketServer()
	logging.Info("WebSocket server started")
}

func (a *Agent) Close() {
	a.DBConn.Close()
}
