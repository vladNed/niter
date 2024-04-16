package main

import (
	"net/http"

	"github.com/indexone/signalling-server/internal/logging"
	"github.com/indexone/signalling-server/internal/settings"
	"github.com/indexone/signalling-server/internal/routes"
	"github.com/indexone/signalling-server/internal/hub"
)

var(
	config = settings.GetSettings()
	logger = logging.GetLogger(nil)
)

func main() {
	// Initialize hub
	go hub.HubInstance.Run()

	// Register routes
	mux := http.NewServeMux()
	routes.HttpRouter.Register(mux)
	routes.WSRouter.Register(mux)

	// Start server
	logger.Info("Running server on: ", config.GetAddress())
	err := http.ListenAndServe(config.GetAddress(), mux)
	if err != nil {
		logger.Error("Error starting server: ", err)
	}
}
