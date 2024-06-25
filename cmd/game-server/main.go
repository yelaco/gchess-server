package main

import (
	"github.com/yelaco/robinhood-chess/internal/agent"
	"github.com/yelaco/robinhood-chess/internal/api"
	"github.com/yelaco/robinhood-chess/pkg/config"
	"github.com/yelaco/robinhood-chess/pkg/logging"
	"go.uber.org/zap"
)

func main() {
	agent := agent.NewAgent()

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
