package http

import (
	"context"
	"github.com/zhayt/kmf-tt/config"
	"github.com/zhayt/kmf-tt/transport/handler"
	"net"
	"net/http"
	"time"
)

type Server struct {
	handler *handler.Handler
	srv     *http.Server
	Notify  chan error
}

func NewServer(cfg *config.Config, handler *handler.Handler) *Server {
	srv := &http.Server{
		Addr:         net.JoinHostPort(":", cfg.App.Port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &Server{srv: srv, handler: handler, Notify: make(chan error, 1)}
}

func (s *Server) StartServer() {
	s.srv.Handler = s.InitRoute()

	go func() {
		s.Notify <- s.srv.ListenAndServe()
	}()
}

func (s *Server) ShutDown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.srv.Shutdown(ctx)
}
