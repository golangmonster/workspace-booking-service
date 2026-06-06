package user

import (
	"context"
	"fmt"

	"github.com/golangmonster/workspace-booking-service/internal/model/user"
)

func (s *service) GetUserByID(ctx context.Context, userID int64) (*user.User, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return user, nil
}
