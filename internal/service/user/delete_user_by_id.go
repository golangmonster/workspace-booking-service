package user

import (
	"context"
	"fmt"
)

func (s *service) DeleteUserByID(ctx context.Context, userID int64) error {
	err := s.repo.DeleteUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}
