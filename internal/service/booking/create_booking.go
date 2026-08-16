package booking

import (
	"context"
	"fmt"

	"github.com/golangmonster/workspace-booking-service/internal/model/booking"
	"github.com/golangmonster/workspace-booking-service/internal/model/workspace"
)

func (s *service) CreateBooking(ctx context.Context, req CreateBookingRequest) (int64, error) {
	if !req.StartAt.Before(req.EndAt) {
		return 0, booking.ErrInvalidTimeRange
	}

	var id int64

	err := s.repo.InTx(ctx, func(ctx context.Context) error {
		_, err := s.userRepo.GetUserByID(ctx, req.UserID)
		if err != nil {
			return fmt.Errorf("get user by id: %w", err)
		}

		ws, err := s.repo.GetWorkspaceForUpdate(ctx, req.WorkspaceID)
		if err != nil {
			return err
		}

		if ws.Status != workspace.StatusAvailable {
			return workspace.ErrWorkspaceNotAvailable
		}

		existing, err := s.repo.ListBookings(ctx, ListBookingsRequest{
			Filter: &Filter{
				WorkspaceID: &req.WorkspaceID,
				Statuses:    []booking.Status{booking.StatusActive},
			},
		})
		if err != nil {
			return fmt.Errorf("list bookings: %w", err)
		}

		for _, b := range existing.Bookings {
			if req.StartAt.Before(b.EndAt) && b.StartAt.Before(req.EndAt) {
				return booking.ErrBookingOverlap
			}
		}

		id, err = s.repo.CreateBooking(ctx, req)

		return err
	})
	if err != nil {
		return 0, err
	}

	return id, nil
}
