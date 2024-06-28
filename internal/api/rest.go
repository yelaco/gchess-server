package api

import (
	"net/http"

	"github.com/yelaco/go-chess-server/pkg/config"
	"github.com/yelaco/go-chess-server/pkg/logging"
	"go.uber.org/zap"
)

// Start REST server
func StartRESTServer(port string) error {
	http.HandleFunc("POST /api/users", handlerUsersCreate)
	http.HandleFunc("POST /api/login", handlerLogin)
	http.HandleFunc("GET /api/sessions", handlerSessionGet)
	http.HandleFunc("GET /api/sessions/{sessionid}", handlerSessionGetFromID)
	logging.Info("rest server started", zap.String("port", config.RESTPort))

	return http.ListenAndServe(":"+port, nil)
}
