package booking

import (
	"time"

	"github.com/golangmonster/workspace-booking-service/internal/model/booking"
	"github.com/golangmonster/workspace-booking-service/internal/model/page"
)

type CreateBookingRequest struct {
	UserID      int64
	WorkspaceID int64
	StartAt     time.Time
	EndAt       time.Time
	Status      booking.Status
}

type UpdateBookingStatus struct {
	BookingID int64
	Status    booking.Status
}

type Filter struct {
	UserID      *int64
	WorkspaceID *int64
	Statuses    []booking.Status

	// DateFrom/DateTo are truncated to a date and matched against bookings
	// whose [StartAt, EndAt) interval overlaps that range of days.
	DateFrom *time.Time
	DateTo   *time.Time
}

type ListBookingsRequest struct {
	Page   *page.Page
	Filter *Filter
}

type ListBookingsResponse struct {
	Bookings   []*booking.Booking
	TotalCount uint32
}
