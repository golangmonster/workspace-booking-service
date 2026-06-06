package user

import (
	"context"
	"errors"

	model "github.com/golangmonster/workspace-booking-service/internal/model/user"
	dto "github.com/golangmonster/workspace-booking-service/internal/service/user"
	pb "github.com/golangmonster/workspace-booking-service/pkg/api/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	err := i.userService.UpdateUser(ctx, &dto.UpdateUserRequest{
		ID:       req.GetId(),
		Login:    req.GetLogin(),
		FullName: req.GetFullName(),
		Phone:    req.GetPhone(),
	})
	if err != nil {
		if errors.Is(err, model.ErrLoginExists) || errors.Is(err, model.ErrPhoneExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		if errors.Is(err, model.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, err
	}

	return &pb.UpdateUserResponse{}, nil
}
