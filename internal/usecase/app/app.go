package app

import (
	"app/config"
	"app/internal/controller"
	"app/pkg/httpserver"
	"app/pkg/logger"
	"app/pkg/postgres"
	"time"
)

func Run() {
	for i := 0; i < config.GetInt("CONN_ATTEMPTS"); i++ {
		if postgres.GetConn() != nil {
			break
		}
		logger.Err("Unable to connect to Postgres. Retrying...", nil)
		time.Sleep(5 * time.Second)
	}

	controller := controller.NewController()
	controller.GetRoutes()

	handler := httpserver.NewHTTPHandler(controller)

	httpServer := httpserver.NewHTTPRestServer(":"+config.Get("PORT"), handler)

	err := httpServer.Serve()
	if err != nil {
		logger.Err("Error of Server", err)
	}
}
