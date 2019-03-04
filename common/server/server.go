package server

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

// StartAsync starts server asynchronously on a random port.
func StartAsync(handler http.Handler) (*http.Server, string) {
	gin.SetMode(gin.ReleaseMode)
	srv := &http.Server{
		Handler: handler,
	}

	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	go func() {
		if err := srv.Serve(listener); err != nil {
			if err != http.ErrServerClosed {
				panic(err)
			}
		}
	}()

	return srv, listener.Addr().String()
}
