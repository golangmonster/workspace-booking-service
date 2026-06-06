package user

import (
	"context"

	"github.com/golangmonster/workspace-booking-service/internal/model/user"
)

type userRepository interface {
	GetUserByID(ctx context.Context, id int64) (*user.User, error)
	GetUserByLogin(ctx context.Context, login string) (*user.User, error)
	InsertUser(ctx context.Context, user *CreateUserRequest) (*user.User, error)
	DeleteUserByID(ctx context.Context, id int64) error
	UpdateUser(ctx context.Context, user *UpdateUserRequest) error
}
