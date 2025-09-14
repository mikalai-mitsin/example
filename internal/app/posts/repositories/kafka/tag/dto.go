package repositories

import (
	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/tag"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func decodeTag(tag entities.Tag) *examplepb.Tag {
	response := &examplepb.Tag{
		Id:        tag.ID.String(),
		CreatedAt: timestamppb.New(tag.CreatedAt),
		UpdatedAt: timestamppb.New(tag.UpdatedAt),
		DeletedAt: nil,
		PostId:    tag.PostId.String(),
		Value:     tag.Value,
	}
	if tag.DeletedAt != nil {
		response.DeletedAt = timestamppb.New(*tag.DeletedAt)
	}
	if tag.DeletedAt != nil {
		response.DeletedAt = timestamppb.New(*tag.DeletedAt)
	}
	if tag.DeletedAt != nil {
		response.DeletedAt = timestamppb.New(*tag.DeletedAt)
	}
	if tag.DeletedAt != nil {
		response.DeletedAt = timestamppb.New(*tag.DeletedAt)
	}
	return response
}
