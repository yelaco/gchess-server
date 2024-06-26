package api

import (
	"fmt"
	"net/http"

	"github.com/yelaco/go-chess-server/pkg/config"
	"github.com/yelaco/go-chess-server/pkg/logging"
	"go.uber.org/zap"
)

func handleTest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("TEST")
}
func StartRESTServer(port string) error {
	http.HandleFunc("/api/test", handleTest)
	logging.Info("rest server started", zap.String("port", config.RESTPort))

	return http.ListenAndServe(":"+port, nil)
}
