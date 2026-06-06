package user

import (
	"context"
	"errors"

	model "github.com/golangmonster/workspace-booking-service/internal/model/user"
	pb "github.com/golangmonster/workspace-booking-service/pkg/api/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) DeleteUserById(ctx context.Context, req *pb.DeleteUserByIdRequest) (*pb.DeleteUserByIdResponse, error) {
	err := i.userService.DeleteUserByID(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, err
	}

	return &pb.DeleteUserByIdResponse{}, nil
}
