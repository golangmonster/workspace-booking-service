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

func (i *Implementation) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user, err := i.userService.CreateUser(ctx, &dto.CreateUserRequest{
		Login:    req.GetLogin(),
		FullName: req.GetFullName(),
		Phone:    req.GetPhone(),
	})
	if err != nil {
		if errors.Is(err, model.ErrLoginExists) || errors.Is(err, model.ErrPhoneExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		return nil, err
	}

	return &pb.CreateUserResponse{
		Id: user.ID,
	}, nil
}
