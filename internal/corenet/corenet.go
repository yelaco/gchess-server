package corenet

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/yelaco/go-chess-server/pkg/config"
	"github.com/yelaco/go-chess-server/pkg/logging"
	"go.uber.org/zap"
)

type WebSocketServer struct {
	address              string
	upgrader             websocket.Upgrader
	messageHandler       func(*websocket.Conn, *Message, *string)
	connCloseGameHandler func(string)
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

func (s *WebSocketServer) SetMessageHandler(msgHandler func(*websocket.Conn, *Message, *string)) {
	s.messageHandler = msgHandler
}

func (s *WebSocketServer) SetConnCloseGameHandler(ccgHandler func(string)) {
	s.connCloseGameHandler = ccgHandler
}

func (s *WebSocketServer) Start() error {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := s.upgrader.Upgrade(w, r, nil)
		if err != nil {
			logging.Error("failed to upgrade connection", zap.String("error", err.Error()))
			return
		}
		defer conn.Close()
		var connID string
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseMessage, websocket.CloseMessage) {
					logging.Info("connection closed", zap.String("remote_address", conn.RemoteAddr().String()))
				} else if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					logging.Info("unexpected close error", zap.String("remote_address", conn.RemoteAddr().String()))
				} else {
					logging.Info("ws message read error", zap.String("remote_address", conn.RemoteAddr().String()))
				}
				s.connCloseGameHandler(connID)
				break
			}

			msg := Message{}
			if err := json.Unmarshal(message, &msg); err != nil {
				conn.Close()
			}
			s.messageHandler(conn, &msg, &connID)
		}
	})
	logging.Info("websocket server started", zap.String("port", config.Port))
	return http.ListenAndServe(s.address, nil)
}
