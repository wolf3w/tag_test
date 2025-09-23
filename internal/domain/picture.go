package domain

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/wolf3w/tag_test/pkg/pb"
)

type PictureInfo struct {
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (pi PictureInfo) ToGRPC() pb.PictureFile {
	return pb.PictureFile{
		Name:      pi.Name,
		CreatedAt: timestamppb.New(pi.CreatedAt),
		UpdatedAt: timestamppb.New(pi.UpdatedAt),
	}
}
