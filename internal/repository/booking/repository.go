package booking

import "github.com/golangmonster/pgxtransactor"

const (
	uniqueViolationCode = "23505"

	bookingUniqueUserWorkspaceUniqueConstraint = "booking_unique_user_workspace_idx"
)

type repository struct {
	pgxtransactor.Transactor

	pool *pgxtransactor.Pool
}

func New(pool *pgxtransactor.Pool) *repository {
	return &repository{
		pool:       pool,
		Transactor: pool,
	}
}
