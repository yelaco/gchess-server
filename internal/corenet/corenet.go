package corenet

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/yelaco/robinhood-chess/pkg/config"
	"github.com/yelaco/robinhood-chess/pkg/logging"
	"go.uber.org/zap"
)

type WebSocketServer struct {
	address        string
	upgrader       websocket.Upgrader
	messageHandler func(*websocket.Conn, *Message)
}

type Message struct {
	Action string                 `json:"action"`
	Data   map[string]interface{} `json:"data"`
}

func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		address: config.Host + ":" + config.Port,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins
			},
		},
	}
}

func (s *WebSocketServer) SetMessageHandler(msgHandler func(*websocket.Conn, *Message)) {
	s.messageHandler = msgHandler
}

func (s *WebSocketServer) Start() error {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := s.upgrader.Upgrade(w, r, nil)
		if err != nil {
			logging.Error("failed to upgrade connection", zap.String("error", err.Error()))
			return
		}
		defer conn.Close()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseMessage) {
					logging.Info("connection closed", zap.String("remote_address", conn.RemoteAddr().String()))
				} else if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					logging.Info("unexpected close error", zap.String("remote_address", conn.RemoteAddr().String()))
				} else {
					logging.Error("ws message read error", zap.String("remote_address", conn.RemoteAddr().String()))
				}
				break
			}

			msg := Message{}
			json.Unmarshal(message, &msg)
			s.messageHandler(conn, &msg)
		}
	})
	logging.Info("websocket server started", zap.String("port", config.Port))
	return http.ListenAndServe(s.address, nil)
}
