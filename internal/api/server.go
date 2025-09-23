package api

import (
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/wolf3w/tag_test/internal/domain"
	"github.com/wolf3w/tag_test/internal/repo"
	"github.com/wolf3w/tag_test/internal/service"
	"github.com/wolf3w/tag_test/pkg/pb"
)

func RunServer(logger *zap.Logger, config *domain.Config) error {
	address := fmt.Sprintf("%s:%s", config.ServerHost, config.ServerPort)

	fsRepo, err := repo.NewFileStorage(config.RootDir)
	if err != nil {
		return fmt.Errorf("cannot create file storage repository: %w", err)
	}

	srv := grpc.NewServer()
	reflection.Register(srv)

	pictureService := service.NewPictureService(logger, fsRepo)
	pb.RegisterPictureStorageServiceServer(srv, pictureService)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("TCP listen: %w", err)
	}

	logger.Info("Starting GRPC service on " + address)
	if err := srv.Serve(listener); err != nil {
		return fmt.Errorf("GRPC serve: %w", err)
	}

	return nil
}
