package user

import (
	"context"
	"fmt"

	"github.com/golangmonster/workspace-booking-service/internal/model/user"
)

func (s *service) CreateUser(ctx context.Context, req *CreateUserRequest) (*user.User, error) {
	user, err := s.repo.InsertUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("insert user: %w", err)
	}

	return user, nil
}
