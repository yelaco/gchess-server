package corenet

import (
	"log"
	"net/url"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/yelaco/robinhood-chess/pkg/config"
)

var ch = make(chan bool)

func TestWebSocketServer(t *testing.T) {
	go setupWebSocketServer()

	u := url.URL{Scheme: "ws", Host: "localhost:" + config.Port, Path: "/ws"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	tests := []struct {
		name  string
		input Message
		want  string
	}{
		{
			"Matching request",
			Message{
				Action: "matching",
				Data: map[string]interface{}{
					"player_id": "42",
				},
			},
			"matching",
		},
		{
			"Move request",
			Message{
				Action: "move",
				Data: map[string]interface{}{
					"session_id": "42",
					"player_id":  "42",
					"move":       "e2-e4",
				},
			},
			"move",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c.WriteJSON(tt.input)
			var ans Message
			if err := c.ReadJSON(&ans); err != nil {
				t.Error(err)
			}
			if ans.Action != tt.want {
				t.Errorf("got %s, want %s", ans.Action, tt.want)
			} else {
				t.Log(ans.Data)
			}
		})
	}

	if err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
		t.Error(err)
	}

	if closed := <-ch; closed {
		t.Log("conn closed is handled")
	}
}

func setupWebSocketServer() {
	wsServer := NewWebSocketServer()
	wsServer.SetMessageHandler(messageHandler)
	wsServer.SetConnCloseGameHandler(connCloseGameHandler)
	log.Fatal(wsServer.Start())
}

func messageHandler(conn *websocket.Conn, message *Message) {
	type errorResponse struct {
		Type  string `json:"type"`
		Error string `json:"error"`
	}
	switch message.Action {
	case "matching":
		_, ok := message.Data["player_id"].(string)
		if ok {
			conn.WriteJSON(*message)
		} else {
			conn.WriteJSON(errorResponse{
				Type:  "error",
				Error: "invalid data",
			})
		}
	case "move":
		_, playerOK := message.Data["player_id"].(string)
		_, sessionOK := message.Data["session_id"].(string)
		_, moveOK := message.Data["move"].(string)
		if playerOK && sessionOK && moveOK {
			conn.WriteJSON(*message)
		} else {
			conn.WriteJSON(errorResponse{
				Type:  "error",
				Error: "invalid data",
			})
		}
	default:
		conn.WriteJSON(errorResponse{
			Type:  "error",
			Error: "invalid action",
		})
	}
}

func connCloseGameHandler(conn *websocket.Conn) {
	ch <- true
}
