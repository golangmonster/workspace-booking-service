package workspace_booking

import (
	"fmt"

	model "github.com/golangmonster/workspace-booking-service/internal/model/booking"
	kafkapb "github.com/golangmonster/workspace-booking-service/pkg/api/kafka/v1"
	modelpb "github.com/golangmonster/workspace-booking-service/pkg/api/model/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	bookingStatusToProto = map[model.Status]modelpb.BookingStatus{
		model.StatusUnspecified: modelpb.BookingStatus_BOOKING_STATUS_UNSPECIFIED,
		model.StatusActive:      modelpb.BookingStatus_BOOKING_STATUS_ACTIVE,
		model.StatusCompleted:   modelpb.BookingStatus_BOOKING_STATUS_COMPLETED,
		model.StatusCancelled:   modelpb.BookingStatus_BOOKING_STATUS_CANCELLED,
	}
)

func MarshalWorkspaceBooking(booking *model.Booking) ([]byte, error) {
	msg := toProtoMessage(booking)

	raw, err := protojson.MarshalOptions{UseProtoNames: true}.Marshal(msg)
	if err != nil {
		return raw, fmt.Errorf("failed to marshal proto to json: %w", err)
	}

	return raw, nil
}

func toProtoMessage(booking *model.Booking) *kafkapb.WorkspaceBookingMessage {
	return &kafkapb.WorkspaceBookingMessage{
		BookingId:   booking.ID,
		Status:      bookingStatusToProto[booking.Status],
		WorkspaceId: booking.WorkspaceID,
		UserId:      booking.UserID,
		StartAt:     timestamppb.New(booking.StartAt),
		EndAt:       timestamppb.New(booking.EndAt),
	}
}
