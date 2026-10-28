package workspace

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	model "github.com/golangmonster/workspace-booking-service/internal/model/workspace"
	pb "github.com/golangmonster/workspace-booking-service/pkg/api/workspace/v1"
)

func (i *Implementation) DeleteWorkspaceById(ctx context.Context, req *pb.DeleteWorkspaceByIdRequest) (*pb.DeleteWorkspaceByIdResponse, error) {
	err := i.workspaceService.DeleteWorkspaceByID(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, model.ErrWorkspaceNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, err
	}

	return &pb.DeleteWorkspaceByIdResponse{}, nil
}
