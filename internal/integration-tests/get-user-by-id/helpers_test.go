package getuserbyid

import (
	"github.com/golangmonster/pgxtransactor"
	userApp "github.com/golangmonster/workspace-booking-service/internal/app/user"
	userRepository "github.com/golangmonster/workspace-booking-service/internal/repository/user"
	userService "github.com/golangmonster/workspace-booking-service/internal/service/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Build the same user dependency chain as main, without starting network servers.
func newUserController(pool *pgxpool.Pool) *userApp.Implementation {
	repository := userRepository.New(pgxtransactor.New(pool))
	service := userService.New(repository)

	return userApp.New(service)
}
