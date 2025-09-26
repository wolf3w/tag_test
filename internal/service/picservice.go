package service

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/wolf3w/tag_test/internal/repo"
	"github.com/wolf3w/tag_test/pkg/pb"
	"github.com/wolf3w/tag_test/pkg/workerpool"
)

const (
	MaxLoadTasks = 10
	MaxWorkers   = 10
	MaxListTasks = 100
)

type PictureService struct {
	pb.UnimplementedPictureStorageServiceServer
	log      *zap.Logger
	repo     *repo.FileStorage
	loadPool workerpool.WorkerPool
	listPool workerpool.WorkerPool
}

func NewPictureService(logger *zap.Logger, storageRepository *repo.FileStorage) *PictureService {
	return &PictureService{
		log:  logger,
		repo: storageRepository,
	}
}

func (sr *PictureService) UploadPicture(_ context.Context, req *pb.PictureUploadRequest) (*emptypb.Empty, error) {
	fileName := req.GetName()
	receivedData := req.GetData()

	err := sr.repo.Write(fileName, receivedData)
	if err != nil {
		sr.log.Error("Write picture", zap.Error(err))
		return nil, fmt.Errorf("write into repo: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (sr *PictureService) ListStoredPictures(_ context.Context, _ *emptypb.Empty) (*pb.ListPicturesResponse, error) {
	rawInfo, err := sr.repo.ListPictures()
	if err != nil {
		sr.log.Error("List pictures", zap.Error(err))
		return nil, fmt.Errorf("list pictures: %w", err)
	}

	picInfo := make([]*pb.PictureFile, len(rawInfo))
	for i := range rawInfo {
		picInfo[i] = &pb.PictureFile{
			Name:      rawInfo[i].Name,
			CreatedAt: timestamppb.New(rawInfo[i].CreatedAt),
			UpdatedAt: timestamppb.New(rawInfo[i].UpdatedAt),
		}
	}

	return &pb.ListPicturesResponse{Pictures: picInfo}, nil
}

func (sr *PictureService) DownloadPicture(
	_ context.Context,
	req *pb.DownloadPictureRequest,
) (*pb.DownloadPictureResponse, error) {
	fileName := req.GetFileName()
	picData, err := sr.repo.Read(fileName)
	if err != nil {
		sr.log.Error("Read picture", zap.Error(err))
		return nil, fmt.Errorf("read picture: %w", err)
	}

	resp := &pb.DownloadPictureResponse{
		Data: picData,
	}
	return resp, nil
}
