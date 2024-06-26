package main

import (
	"github.com/yelaco/go-chess-server/internal/agent"
	"github.com/yelaco/go-chess-server/internal/api"
	"github.com/yelaco/go-chess-server/pkg/config"
	"github.com/yelaco/go-chess-server/pkg/logging"
	"go.uber.org/zap"
)

func main() {
	agent := agent.NewAgent()

	// log.Fatal(agent.StartGameServer())

	go func() {
		if err := agent.StartGameServer(); err != nil {
			logging.Fatal("game server failed to start", zap.Error(err))
		}
	}()

	go func() {
		if err := api.StartRESTServer(config.RESTPort); err != nil {
			logging.Fatal("rest server failed to start", zap.Error(err))
		}
	}()

	select {}
}
