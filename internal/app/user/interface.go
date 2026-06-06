package user

import (
	"context"

	model "github.com/golangmonster/workspace-booking-service/internal/model/user"
	dto "github.com/golangmonster/workspace-booking-service/internal/service/user"
)

type userService interface {
	GetUserByID(ctx context.Context, userID int64) (*model.User, error)
	CreateUser(ctx context.Context, user *dto.CreateUserRequest) (*model.User, error)
	UpdateUser(ctx context.Context, user *dto.UpdateUserRequest) error
	DeleteUserByID(ctx context.Context, userID int64) error
}
