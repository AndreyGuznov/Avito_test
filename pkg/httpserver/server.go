package httpserver

import (
	"app/pkg/logger"
	"net/http"
)

// HTTPRestServer ...
type HTTPRestServer struct {
	address string
	wrapped *http.Server
}

// NewHTTPRestServer creates HTTP service
func NewHTTPRestServer(address string, handler http.Handler) *HTTPRestServer {
	httpsrv := http.Server{
		Addr:    address,
		Handler: handler,
	}

	return &HTTPRestServer{wrapped: &httpsrv, address: address}
}

// Serve HTTP service

func (server *HTTPRestServer) Serve() error {
	logger.Info("Serving on " + server.address)
	err := server.wrapped.ListenAndServe()

	if err != http.ErrServerClosed {
		logger.Err("Server crashed", err)
	}

	return err
}
