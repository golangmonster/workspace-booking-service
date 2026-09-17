package booking

import (
	"context"

	"github.com/golangmonster/workspace-booking-service/internal/model/booking"
	"github.com/golangmonster/workspace-booking-service/internal/model/outbox"
	"github.com/golangmonster/workspace-booking-service/internal/model/user"
	"github.com/golangmonster/workspace-booking-service/internal/model/workspace"
)

type bookingRepository interface {
	CreateBooking(ctx context.Context, req CreateBookingRequest) (*booking.Booking, error)
	UpdateBookingStatus(ctx context.Context, req UpdateBookingStatus) (*booking.Booking, error)
	ListBookings(ctx context.Context, req ListBookingsRequest) (ListBookingsResponse, error)
	GetWorkspaceForUpdate(ctx context.Context, workspaceID int64) (*workspace.Workspace, error)
	InTx(ctx context.Context, fn func(ctx context.Context) error) error
	GetBookingWithLock(ctx context.Context, bookingID int64) (*booking.Booking, error)
}

type userRepository interface {
	GetUserByID(ctx context.Context, userID int64) (*user.User, error)
}

type outboxRepository interface {
	CreateOutboxItem(ctx context.Context, item outbox.Message) error
}
