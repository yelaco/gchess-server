package main

import (
	"github.com/yelaco/robinhood-chess/internal/agent"
	"github.com/yelaco/robinhood-chess/internal/matcher"
)

func main() {
	ag := agent.NewAgent()
	ag.StartSocketServer()

	matcher := matcher.NewMatcher()

}
