//go:build integration

package deleteuserbyid_test

import (
	"testing"

	"github.com/golangmonster/pgxtransactor"

	userApp "github.com/golangmonster/workspace-booking-service/internal/app/user"
	"github.com/golangmonster/workspace-booking-service/internal/integration-test/testcase"
	"github.com/golangmonster/workspace-booking-service/internal/integration-test/testdb"
	userRepository "github.com/golangmonster/workspace-booking-service/internal/repository/user"
	userService "github.com/golangmonster/workspace-booking-service/internal/service/user"
	pb "github.com/golangmonster/workspace-booking-service/pkg/api/user/v1"
)

func TestDeleteUserByID(t *testing.T) {
	pg := testdb.Start(t, "migrations")

	for _, tc := range testcase.Load(t, "scenario.yaml") {
		t.Run(tc.Name, func(t *testing.T) {
			t.Log(tc.Description)

			db := pg.NewDatabase(t, tc.Path("init.yaml"))

			repo := userRepository.New(pgxtransactor.New(db.Pool))
			controller := userApp.New(userService.New(repo))

			var request pb.DeleteUserByIdRequest
			tc.ReadRequest(t, &request)

			response, err := controller.DeleteUserById(t.Context(), &request)
			tc.AssertResponse(t, response, err)
			db.Assert(t, tc.Path("assert.yaml"))
		})
	}
}
