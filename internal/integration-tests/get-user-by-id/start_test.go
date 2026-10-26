package getuserbyid

import (
	"testing"
	"time"

	"github.com/golangmonster/workspace-booking-service/internal/integration-tests/testdb"
	pb "github.com/golangmonster/workspace-booking-service/pkg/api/user/v1"
	"github.com/stretchr/testify/require"
	"github.com/viant/dsunit"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGetUserByID(t *testing.T) {
	cases := []struct {
		name         string
		prepareDir   string
		expectDir    string
		request      *pb.GetUserByIdRequest
		wantResponse *pb.GetUserByIdResponse
	}{
		{
			name:       "returns existing user",
			prepareDir: "testdata/existing-user/prepare",
			expectDir:  "testdata/existing-user/expect",
			request:    &pb.GetUserByIdRequest{Id: 1},
			wantResponse: &pb.GetUserByIdResponse{
				Id:        1,
				Login:     "user1",
				FullName:  "Test User",
				Phone:     "+79990001122",
				CreatedAt: timestamppb.New(time.Date(2026, time.January, 1, 10, 0, 0, 0, time.UTC)),
				UpdatedAt: timestamppb.New(time.Date(2026, time.January, 1, 11, 0, 0, 0, time.UTC)),
			},
		},
	}

	postgres := testdb.StartPostgres(t)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pool := testdb.NewDatabase(t, postgres)
			fixtures := testdb.NewFixtures(t, pool)
			datastore := pool.Config().ConnConfig.Database
			controller := newUserController(pool)

			initial := dsunit.NewDatasetResource(datastore, tc.prepareDir, "", "")
			prepared := fixtures.Prepare(dsunit.NewPrepareRequest(initial))
			require.NoError(t, prepared.Error(), "prepare database fixtures")

			response, err := controller.GetUserById(t.Context(), tc.request)

			require.NoError(t, err)
			require.True(t, proto.Equal(tc.wantResponse, response), "response mismatch:\nwant: %s\ngot: %s",
				protojson.Format(tc.wantResponse), protojson.Format(response))

			expected := dsunit.NewDatasetResource(datastore, tc.expectDir, "", "")
			checked := fixtures.Expect(dsunit.NewExpectRequest(dsunit.FullTableDatasetCheckPolicy, expected))
			require.NoError(t, checked.Error(), "database must match expected fixtures")
		})
	}
}
