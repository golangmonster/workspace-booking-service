package workspace

import (
	"github.com/golangmonster/pgxtransactor"
)

const (
	uniqueViolationCode = "23505"

	workspaceNameTypeUniqueConstraint = "workspace_unique_name_type_idx"
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
