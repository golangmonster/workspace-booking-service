package user

import (
	"context"
	"fmt"
)

func (s *service) UpdateUser(ctx context.Context, req *UpdateUserRequest) error {
	err := s.repo.UpdateUser(ctx, req)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	return nil
}
