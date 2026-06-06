package user

import (
	"github.com/golangmonster/pgxtransactor"
)

const (
	uniqueViolationCode = "23505"

	loginUniqueConstraintName = "unique_login"
	phoneUniqueConstraintName = "unique_phone"
)

type repository struct {
	pool *pgxtransactor.Pool
	pgxtransactor.Transactor
}

func New(pool *pgxtransactor.Pool) *repository {
	return &repository{
		pool:       pool,
		Transactor: pool,
	}
}
