package main

import (
	"fmt"
	"github.com/zhayt/kmf-tt/config"
	"github.com/zhayt/kmf-tt/logger"
	"github.com/zhayt/kmf-tt/service"
	"github.com/zhayt/kmf-tt/storage"
	"github.com/zhayt/kmf-tt/transport/handler"
	"github.com/zhayt/kmf-tt/transport/http"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	l, err := logger.Init(cfg)
	if err != nil {
		return fmt.Errorf("cannot init logger: %w", err)
	}
	defer func(l *zap.Logger) {
		err = l.Sync()
		if err != nil {
			log.Fatalln(err)
		}
	}(l)

	// repo
	repo, err := storage.NewStorage(cfg, l)
	if err != nil {
		return err
	}

	// service
	servi := service.NewCurrencyService(repo, l)

	// handler
	hand := handler.NewHandler(servi, l)

	// server
	srv := http.NewServer(cfg, hand)

	l.Info("Start app", zap.String("port", cfg.App.Port))
	srv.StartServer()

	// grace full shutdown
	osSignCh := make(chan os.Signal, 1)
	signal.Notify(osSignCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-osSignCh:
		l.Info("signal accepted: ", zap.String("signal", s.String()))
	case err = <-srv.Notify:
		l.Info("server closing", zap.Error(err))
	}

	if err = srv.ShutDown(); err != nil {
		return fmt.Errorf("error while shutting down server: %s", err)
	}
	return nil
}
