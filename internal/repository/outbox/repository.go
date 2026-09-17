package outbox

import "github.com/golangmonster/pgxtransactor"

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
