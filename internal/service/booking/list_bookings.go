package booking

import (
	"context"
	"fmt"
)

func (s *service) ListBookings(ctx context.Context, req ListBookingsRequest) (*ListBookingsResponse, error) {
	resp, err := s.repo.ListBookings(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("listing bookings: %w", err)
	}

	return resp, nil
}
