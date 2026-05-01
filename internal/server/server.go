package server

import (
	"context"
	"fmt"
	"net"

	"github.com/aclgo/grpc-admin/config"
	"github.com/aclgo/grpc-admin/internal/admin"
	"github.com/aclgo/grpc-admin/internal/admin/repository"
	"github.com/aclgo/grpc-admin/internal/admin/usecase"
	"github.com/aclgo/grpc-admin/internal/delivery/grpc/service"
	"github.com/aclgo/grpc-admin/pkg/logger"
	proto "github.com/aclgo/grpc-admin/proto/admin"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

type Server struct {
	config        *config.Config
	db            *sqlx.DB
	logger        logger.Logger
	observability *admin.Observability
	redis         *redis.Client
}

func NewServer(cfg *config.Config, db *sqlx.DB, rds *redis.Client, logger logger.Logger, obs *admin.Observability) *Server {
	return &Server{
		config:        cfg,
		db:            db,
		redis:         rds,
		logger:        logger,
		observability: obs,
	}
}

func (s *Server) Run(ctx context.Context) error {
	adminRepo := repository.NewpostgresRepo(s.db)
	adminUC := usecase.NewAdminService(adminRepo, s.redis, s.logger)

	handlers := service.NewAdminService(adminUC, s.observability)

	// ctx := context.Background()
	interceptor := NewInterceptor(s.logger)

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(interceptor.GrpcInterceptor),
		// grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		// grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
	}

	srv := grpc.NewServer(opts...)

	proto.RegisterAdminServiceServer(srv, handlers)

	listen, err := net.Listen("tcp", fmt.Sprintf(":%v", s.config.ApiPort))
	if err != nil {
		return errors.Wrap(err, "Run.Listen")
	}

	ec := make(chan error)

	go func() {
		s.logger.Infof("server starting port %v", s.config.ApiPort)
		ec <- srv.Serve(listen)
	}()

	select {
	case <-ec:
		if err != nil {
			return errors.Wrap(<-ec, "Run.Serve")
		}

	case <-ctx.Done():
		s.logger.Info("server stop")
		srv.GracefulStop()
	}

	return nil
}
