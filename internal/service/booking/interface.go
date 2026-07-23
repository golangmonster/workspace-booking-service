package booking

import (
	"context"

	"github.com/golangmonster/workspace-booking-service/internal/model/workspace"
)

type bookingRepository interface {
	CreateBooking(ctx context.Context, req CreateBookingRequest) (int64, error)
	UpdateBookingStatus(ctx context.Context, req UpdateBookingStatus) error
	ListBookings(ctx context.Context, req ListBookingsRequest) (ListBookingsResponse, error)
	GetWorkspaceForUpdate(ctx context.Context, workspaceID int64) (*workspace.Workspace, error)
	InTx(ctx context.Context, fn func(ctx context.Context) error) error
}
