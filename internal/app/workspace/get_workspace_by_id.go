package workspace

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	model "github.com/golangmonster/workspace-booking-service/internal/model/workspace"
	pb "github.com/golangmonster/workspace-booking-service/pkg/api/workspace/v1"
)

func (i *Implementation) GetWorkspaceById(ctx context.Context, req *pb.GetWorkspaceByIdRequest) (*pb.GetWorkspaceByIdResponse, error) {
	workspace, err := i.workspaceService.GetWorkspaceByID(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, model.ErrWorkspaceNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, err
	}

	return &pb.GetWorkspaceByIdResponse{
		Workspace: &pb.GetWorkspaceByIdResponse_Workspace{
			Id:          workspace.ID,
			Name:        workspace.Name,
			Coordinates: toLatLon(workspace.Lat, workspace.Lon),
			FullAddress: workspace.FullAddress,
			Type:        workspaceTypeToProto[workspace.Type],
			Status:      workspaceStatusToProto[workspace.Status],
			Capacity:    workspace.Capacity,
			CreatedAt:   timestamppb.New(workspace.CreatedAt),
			UpdatedAt:   timestamppb.New(workspace.UpdatedAt),
		},
	}, nil
}
