package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/aclgo/grpc-admin/config"
	"github.com/aclgo/grpc-admin/internal/server"
	"github.com/aclgo/grpc-admin/pkg/logger"
	"github.com/aclgo/grpc-admin/pkg/postgres"
	"github.com/aclgo/grpc-admin/pkg/rds"
)

func main() {

	cfg := config.NewConfig(".")

	logger := logger.NewapiLogger(cfg)

	// otel, err := tel.NewOtel(cfg, logger, "grpc-admin", "0.01")
	// if err != nil {
	// 	logger.Errorf("otel.NewOtel: %v", err)
	// 	return
	// }

	// defer func() {
	// 	otel.Shutdowns(context.Background())
	// }()

	// mt := otel.MeterProvider.Meter("grpc-jwt")
	// tr := otel.TracerProvider.Tracer("grpc-jwt")

	db := postgres.Connect(cfg.DbDriver, cfg.DbURI)
	redis := rds.Connect(cfg)

	server := server.NewServer(cfg, db, redis, logger, nil)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// go func() {
	// 	time.Sleep(time.Second * 10)
	// 	e2e.Run(cfg)
	// }()

	if err := server.Run(ctx); err != nil {
		logger.Errorf("server.Run: %v", err)
	}
}
