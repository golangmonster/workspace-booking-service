//go:build integration

package cancelbooking_test

import (
	"testing"

	"github.com/golangmonster/pgxtransactor"

	bookingApp "github.com/golangmonster/workspace-booking-service/internal/app/booking"
	"github.com/golangmonster/workspace-booking-service/internal/integration-test/testcase"
	"github.com/golangmonster/workspace-booking-service/internal/integration-test/testdb"
	workspaceBooking "github.com/golangmonster/workspace-booking-service/internal/kafka/producer/workspace-booking"
	bookingRepository "github.com/golangmonster/workspace-booking-service/internal/repository/booking"
	outboxRepository "github.com/golangmonster/workspace-booking-service/internal/repository/outbox"
	bookingService "github.com/golangmonster/workspace-booking-service/internal/service/booking"
	pb "github.com/golangmonster/workspace-booking-service/pkg/api/booking/v1"
)

func TestCancelBooking(t *testing.T) {
	pg := testdb.Start(t, "migrations")

	for _, tc := range testcase.Load(t, "scenario.yaml") {
		t.Run(tc.Name, func(t *testing.T) {
			t.Log(tc.Description)

			db := pg.NewDatabase(t, tc.Path("init.yaml"))

			tx := pgxtransactor.New(db.Pool)
			srv := bookingService.New(
				bookingRepository.New(tx),
				nil,
				outboxRepository.New(tx),
				workspaceBooking.MarshalWorkspaceBooking,
				"workspace-booking",
			)

			controller := bookingApp.New(srv)

			var request pb.CancelBookingRequest
			tc.ReadRequest(t, &request)

			response, err := controller.CancelBooking(t.Context(), &request)
			tc.AssertResponse(t, response, err)
			db.Assert(t, tc.Path("assert.yaml"))
		})
	}
}
