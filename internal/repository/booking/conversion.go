package booking

import "github.com/golangmonster/workspace-booking-service/internal/model/booking"

func toBooking(i bookingItem) *booking.Booking {
	return &booking.Booking{
		ID:          i.ID,
		WorkspaceID: i.WorkspaceID,
		UserID:      i.UserID,
		StartAt:     i.StartAt,
		EndAt:       i.EndAt,
		Status:      booking.Status(i.Status),
		CreatedAt:   i.CreatedAt,
		UpdatedAt:   i.UpdatedAt,
	}
}
