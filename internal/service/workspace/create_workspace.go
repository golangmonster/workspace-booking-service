package workspace

import (
	"context"
	"fmt"

	model "github.com/golangmonster/workspace-booking-service/internal/model/workspace"
)

func (s *service) CreateWorkspace(ctx context.Context, req *CreateWorkspaceRequest) (*model.Workspace, error) {
	req.Status = model.StatusAvailable

	workspace, err := s.repo.InsertWorkspace(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("insert workspace: %w", err)
	}

	return workspace, nil
}
