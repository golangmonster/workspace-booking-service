package workspace

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	model "github.com/golangmonster/workspace-booking-service/internal/model/workspace"
	dto "github.com/golangmonster/workspace-booking-service/internal/service/workspace"
	pb "github.com/golangmonster/workspace-booking-service/pkg/api/workspace/v1"
)

func (i *Implementation) UpdateWorkspace(ctx context.Context, req *pb.UpdateWorkspaceRequest) (*pb.UpdateWorkspaceResponse, error) {
	err := i.workspaceService.UpdateWorkspace(ctx, &dto.UpdateWorkspaceRequest{
		ID:          req.GetId(),
		Name:        req.GetName(),
		Lat:         req.GetCoordinates().GetLat(),
		Lon:         req.GetCoordinates().GetLon(),
		FullAddress: req.GetFullAddress(),
		Status:      workspaceStatusToModel[req.GetStatus()],
		Capacity:    req.GetCapacity(),
	})
	if err != nil {
		if errors.Is(err, model.ErrWorkspaceNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, err
	}

	return &pb.UpdateWorkspaceResponse{}, nil
}
