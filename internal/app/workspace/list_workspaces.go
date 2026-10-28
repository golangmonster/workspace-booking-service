package workspace

import (
	"context"

	pb "github.com/golangmonster/workspace-booking-service/pkg/api/workspace/v1"
)

func (i *Implementation) ListWorkspaces(ctx context.Context, req *pb.ListWorkspacesRequest) (*pb.ListWorkspacesResponse, error) {
	resp, err := i.workspaceService.ListWorkspaces(ctx, toListWorkspacesRequest(req))
	if err != nil {
		return nil, err
	}

	return toListWorkspacesResponse(resp), nil
}
