package grpc

import (
	"github.com/himbo22/xoxz/account-service/internal/model"
	media "github.com/himbo22/xoxz/common-service/protobuf/media"
)

func ToCommitFileRequest(req model.UpdateAvatarRequest) *media.CommitFileRequest {
	return &media.CommitFileRequest{
		TmpPath: req.TmpPath,
		PerPath: req.PerPath,
	}
}
