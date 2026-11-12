package main

import (
	"context"
	"errors"
	"flag"
	stdlog "log"
	"os/signal"
	"syscall"

	"github.com/dvadert/chat-service/internal/config"
	"github.com/dvadert/chat-service/internal/logger"
	serverdebug "github.com/dvadert/chat-service/internal/server-debug"
)

func main() {
	configPath := flag.String("configPath", "configs/config.toml", "path to config file")
	flag.Parse()

	if err := run(*configPath); err != nil {
		stdlog.Fatalf("run: %v", err)
	}
}

func run(configPath string) error {
	cfg, err := config.ParseAndValidate(configPath)
	if err != nil {
		return err
	}

	if err := logger.Init(logger.NewOptions(
		cfg.Log.Level,
		logger.WithProductionMode(cfg.Global.Env == "prod"),
	)); err != nil {
		return err
	}
	defer logger.Sync()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv, err := serverdebug.New(serverdebug.NewOptions(cfg.Servers.Debug.Addr))
	if err != nil {
		return err
	}

	if err := srv.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}
