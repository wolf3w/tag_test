package main_test

import (
	"context"
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/wolf3w/tag_test/internal/app"
	"github.com/wolf3w/tag_test/internal/domain"
	"github.com/wolf3w/tag_test/pkg/pb"
)

var (
	//go:embed testdata/case_1.png
	baseImgOne []byte

	//go:embed testdata/case_2.png
	baseImgTwo []byte

	//go:embed testdata/case_3.png
	baseImgThree []byte
)

type ClientSuite struct {
	suite.Suite
	address string
	conn    *grpc.ClientConn
	client  pb.PictureStorageServiceClient
}

func TestClient(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}

func (cs *ClientSuite) SetupSuite() {
	loggerStub := zap.Must(zap.NewDevelopment())
	cfg, err := domain.NewFromEnv()
	cs.Require().NoError(err)

	cs.address = fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)

	serverApp := app.NewApp(loggerStub, cfg)

	go func() {
		if err := serverApp.Run(); err != nil {
			cs.Require().NoError(err)
		}
	}()
}

func (cs *ClientSuite) SetupTest() {
	var err error

	cs.conn, err = grpc.NewClient(cs.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	cs.Require().NoError(err)

	cs.client = pb.NewPictureStorageServiceClient(cs.conn)
}

func (cs *ClientSuite) TearDownTest() {
	err := cs.conn.Close()
	cs.Require().NoError(err)
}

func (cs *ClientSuite) TestUpload() {
	ctx := context.Background()

	req := &pb.PictureUploadRequest{
		Name: "test_1.png",
		Data: baseImgOne,
	}

	_, err := cs.client.UploadPicture(ctx, req)
	cs.Require().NoError(err)

	resp, err := cs.client.ListStoredPictures(ctx, nil)
	cs.Require().NoError(err)

	cs.Require().NotNil(resp.Pictures)
	cs.Require().Equal(resp.Pictures[0].Name, "test_1.png")
}

func (cs *ClientSuite) TestUploadAndDownload() {
	ctx := context.Background()

	reqOne := &pb.PictureUploadRequest{
		Name: "test_2.png",
		Data: baseImgTwo,
	}
	reqTwo := &pb.PictureUploadRequest{
		Name: "test_3.png",
		Data: baseImgThree,
	}

	_, err := cs.client.UploadPicture(ctx, reqOne)
	cs.Require().NoError(err)

	_, err = cs.client.UploadPicture(ctx, reqTwo)
	cs.Require().NoError(err)

	listResp, err := cs.client.ListStoredPictures(ctx, nil)
	cs.Require().NoError(err)

	cs.Require().NotNil(listResp.Pictures)
	cs.Require().Len(listResp.Pictures, 3)

	downloadReq := &pb.DownloadPictureRequest{FileName: "test_3.png"}

	downloadResp, err := cs.client.DownloadPicture(ctx, downloadReq)
	cs.Require().NoError(err)

	cs.Require().NotNil(downloadResp)
	cs.Require().EqualValues(baseImgThree, downloadResp.Data)
}
