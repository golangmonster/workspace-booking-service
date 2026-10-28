package workspace

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	model "github.com/golangmonster/workspace-booking-service/internal/model/workspace"
	pb "github.com/golangmonster/workspace-booking-service/pkg/api/workspace/v1"
)

func (i *Implementation) CreateWorkspace(ctx context.Context, req *pb.CreateWorkspaceRequest) (*pb.CreateWorkspaceResponse, error) {
	workspace, err := i.workspaceService.CreateWorkspace(ctx, toCreateWorkspaceRequest(req))
	if err != nil {
		if errors.Is(err, model.ErrWorkspaceAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		return nil, err
	}

	return &pb.CreateWorkspaceResponse{
		WorkspaceId: workspace.ID,
	}, nil
}
