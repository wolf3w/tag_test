package service

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/wolf3w/tag_test/internal/domain"
	"github.com/wolf3w/tag_test/internal/repo"
	"github.com/wolf3w/tag_test/pkg/pb"
)

const (
	MaxLoadTasks = 10
	MaxListTasks = 100
)

type PictureService struct {
	pb.UnimplementedPictureStorageServiceServer
	log    *zap.Logger
	repo   *repo.FileStorage
	loadCh chan struct{}
	listCh chan struct{}
}

func NewPictureService(logger *zap.Logger, storageRepository *repo.FileStorage) *PictureService {
	return &PictureService{
		log:    logger,
		repo:   storageRepository,
		loadCh: make(chan struct{}, MaxLoadTasks),
		listCh: make(chan struct{}, MaxListTasks),
	}
}

func (sr *PictureService) UploadPicture(ctx context.Context, req *pb.PictureUploadRequest) (*emptypb.Empty, error) {
	// TODO: Загуглить делает ли gRPC-сервер каждый обработчик в своём потоке
	sr.loadCh <- struct{}{}
	resultCh := make(chan error)

	fileName := req.GetName()
	receivedData := req.GetData()

	go func() {
		var err error

		err = sr.repo.Write(fileName, receivedData)
		if err != nil {
			sr.log.Error("Write picture", zap.Error(err))
			err = fmt.Errorf("write into repo: %w", err)
		}

		select {
		case resultCh <- err:
		case <-ctx.Done():
		}

		close(resultCh)
		<-sr.loadCh
	}()

	// Ждём пока горутина не запишет результат
	err := <-resultCh
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (sr *PictureService) ListStoredPictures(ctx context.Context, _ *emptypb.Empty) (*pb.ListPicturesResponse, error) {
	sr.listCh <- struct{}{}
	resultCh := make(chan domain.RespPair[*pb.ListPicturesResponse])

	go func() {
		listObj, err := sr.listPictures()

		select {
		case resultCh <- domain.RespPair[*pb.ListPicturesResponse]{Resp: listObj, Err: err}:
		case <-ctx.Done():
		}

		close(resultCh)
		<-sr.listCh
	}()

	res := <-resultCh
	if res.Err != nil {
		return nil, res.Err
	}
	return res.Resp, nil
}

func (sr *PictureService) DownloadPicture(
	ctx context.Context,
	req *pb.DownloadPictureRequest,
) (*pb.DownloadPictureResponse, error) {
	sr.loadCh <- struct{}{}
	resultCh := make(chan domain.RespPair[*pb.DownloadPictureResponse])

	fileName := req.GetFileName()

	go func() {
		resp := domain.RespPair[*pb.DownloadPictureResponse]{
			Resp: &pb.DownloadPictureResponse{},
		}
		picData, err := sr.repo.Read(fileName)
		if err != nil {
			sr.log.Error("Read picture", zap.Error(err))
			err = fmt.Errorf("read picture: %w", err)
			resp.Err = err
		}
		resp.Resp.Data = picData

		select {
		case resultCh <- resp:
		case <-ctx.Done():
		}

		close(resultCh)
		<-sr.loadCh
	}()

	res := <-resultCh
	if res.Err != nil {
		return nil, res.Err
	}
	return res.Resp, nil
}

func (sr *PictureService) listPictures() (*pb.ListPicturesResponse, error) {
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
