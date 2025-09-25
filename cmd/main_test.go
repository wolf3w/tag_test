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

	sender, err := cs.client.UploadPicture(ctx)
	cs.Require().NoError(err)

	req := &pb.PictureUploadRequest{
		Name: "test_1.png",
		Data: baseImgOne,
	}

	err = sender.Send(req)
	cs.Require().NoError(err)

	err = sender.CloseSend()
	cs.Require().NoError(err)

	list, err := cs.client.ListStoredPictures(ctx, nil)
	cs.Require().NoError(err)

	cs.Require().NotNil(list.Pictures)
	cs.Require().Equal(list.Pictures[0].Name, "test_1.png")
}

func (cs *ClientSuite) TestUploadAndDownload() {
	ctx := context.Background()

	sender, err := cs.client.UploadPicture(ctx)
	cs.Require().NoError(err)

	reqOne := &pb.PictureUploadRequest{
		Name: "test_2.png",
		Data: baseImgTwo,
	}
	reqTwo := &pb.PictureUploadRequest{
		Name: "test_3.png",
		Data: baseImgThree,
	}

	err = sender.Send(reqOne)
	cs.Require().NoError(err)

	err = sender.Send(reqTwo)
	cs.Require().NoError(err)

	err = sender.CloseSend()
	cs.Require().NoError(err)

}
