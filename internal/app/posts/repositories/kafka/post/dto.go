package repositories

import (
	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/post"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func decodePost(post entities.Post) *examplepb.Post {
	response := &examplepb.Post{
		Id:        post.ID.String(),
		CreatedAt: timestamppb.New(post.CreatedAt),
		UpdatedAt: timestamppb.New(post.UpdatedAt),
		DeletedAt: nil,
		Body:      post.Body,
	}
	if post.DeletedAt != nil {
		response.DeletedAt = timestamppb.New(*post.DeletedAt)
	}
	return response
}
