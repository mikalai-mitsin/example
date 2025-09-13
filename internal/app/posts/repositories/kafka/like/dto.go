package repositories

import (
	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/like"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func decodeLike(like entities.Like) *examplepb.Like {
	response := &examplepb.Like{
		Id:        like.ID.String(),
		CreatedAt: timestamppb.New(like.CreatedAt),
		UpdatedAt: timestamppb.New(like.UpdatedAt),
		DeletedAt: nil,
		PostId:    like.PostId.String(),
		Value:     like.Value,
		UserId:    like.UserId.String(),
	}
	if like.DeletedAt != nil {
		response.DeletedAt = timestamppb.New(*like.DeletedAt)
	}
	if like.DeletedAt != nil {
		response.DeletedAt = timestamppb.New(*like.DeletedAt)
	}
	return response
}
