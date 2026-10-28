package user

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	model "github.com/golangmonster/workspace-booking-service/internal/model/user"
	pb "github.com/golangmonster/workspace-booking-service/pkg/api/user/v1"
)

func (i *Implementation) GetUserById(ctx context.Context, req *pb.GetUserByIdRequest) (*pb.GetUserByIdResponse, error) {
	user, err := i.userService.GetUserByID(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, err
	}

	return &pb.GetUserByIdResponse{
		Id:        user.ID,
		Login:     user.Login,
		FullName:  user.FullName,
		Phone:     user.Phone,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}
