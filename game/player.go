package game

import "github.com/yelaco/robinhood-chess/pkg/utils"

type player struct {
	playerId    string
	isWhiteSide bool
}

func initPlayers() *[2]player {
	players := [2]player{
		{
			playerId:    utils.GenerateUUID(),
			isWhiteSide: true,
		},
		{
			playerId:    utils.GenerateUUID(),
			isWhiteSide: false,
		},
	}
	return &players
}
