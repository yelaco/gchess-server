package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"os/signal"

	"github.com/gorilla/websocket"
	"github.com/yelaco/go-chess-server/internal/corenet"
	"github.com/yelaco/go-chess-server/internal/session"
	"github.com/yelaco/go-chess-server/pkg/config"
)

type matchResponse struct {
	Type        string              `json:"type"`
	SessionID   string              `json:"session_id"`
	GameState   session.GameState   `json:"game_state"`
	PlayerState session.PlayerState `json:"player_state"`
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "localhost:" + config.Port, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})
	defer close(done)

	playerID := "11224"

	c.WriteJSON(corenet.Message{
		Action: "matching",
		Data: map[string]interface{}{
			"player_id": playerID,
		},
	})

	var resp matchResponse
	if err := c.ReadJSON(&resp); err != nil {
		log.Fatal("ws match", err)
		return
	} else {
		fmt.Println(resp.Type)
	}

	if resp.Type != "matched" {
		log.Fatal("not matched")
		return
	}

	state := resp.GameState
	scanner := bufio.NewScanner(os.Stdin)
	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			return
		default:
			clearScreen()
			printBoard(state.Board)
			if resp.PlayerState.IsWhiteSide == state.IsWhiteTurn {
				fmt.Print("Enter your move (e.g., e2-e4): ")
				scanner.Scan()
				move := scanner.Text()

				c.WriteJSON(corenet.Message{
					Action: "move",
					Data: map[string]interface{}{
						"session_id": resp.SessionID,
						"player_id":  playerID,
						"move":       move,
					},
				})

				err := c.ReadJSON(&state)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				fmt.Print("Wait for your opponent...")
				if err := c.ReadJSON(&state); err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

func printBoard(board [8][8]string) {
	fmt.Println("   a b c d e f g h")
	fmt.Println("  +-----------------+")

	for i := 7; i >= 0; i-- {
		fmt.Printf("%d | ", i+1)
		for j := 0; j < 8; j++ {
			box := board[j][i]
			if box != "" {
				fmt.Print(box + " ")
			} else {
				fmt.Print(". ")
			}
		}
		fmt.Println("|")
	}

	fmt.Println("  +-----------------+")
	fmt.Println()
}

func clearScreen() {
	cmd := exec.Command("clear") // for Windows use "cls"
	cmd.Stdout = os.Stdout
	cmd.Run()
}
