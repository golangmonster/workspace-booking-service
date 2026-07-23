package booking

import "context"

func (s *service) ListBookings(ctx context.Context, req ListBookingsRequest) (ListBookingsResponse, error) {
	return s.repo.ListBookings(ctx, req)
}
