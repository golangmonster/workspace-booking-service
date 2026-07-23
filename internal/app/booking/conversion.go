package booking

import (
	"github.com/golangmonster/workspace-booking-service/internal/model/booking"
	"github.com/golangmonster/workspace-booking-service/internal/model/page"
	dto "github.com/golangmonster/workspace-booking-service/internal/service/booking"
	pb "github.com/golangmonster/workspace-booking-service/pkg/api/booking/v1"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	bookingStatusToModel = map[pb.BookingStatus]booking.Status{
		pb.BookingStatus_BOOKING_STATUS_UNSPECIFIED: booking.StatusUnspecified,
		pb.BookingStatus_BOOKING_STATUS_ACTIVE:      booking.StatusActive,
		pb.BookingStatus_BOOKING_STATUS_CANCELLED:   booking.StatusCancelled,
		pb.BookingStatus_BOOKING_STATUS_COMPLETED:   booking.StatusCompleted,
	}

	bookingStatusToProto = map[booking.Status]pb.BookingStatus{
		booking.StatusUnspecified: pb.BookingStatus_BOOKING_STATUS_UNSPECIFIED,
		booking.StatusActive:      pb.BookingStatus_BOOKING_STATUS_ACTIVE,
		booking.StatusCancelled:   pb.BookingStatus_BOOKING_STATUS_CANCELLED,
		booking.StatusCompleted:   pb.BookingStatus_BOOKING_STATUS_COMPLETED,
	}
)

func toListBookingsRequest(req *pb.ListBookingsRequest) dto.ListBookingsRequest {
	var filter *dto.Filter

	if req.Filter != nil {
		filter = &dto.Filter{
			UserID:      req.Filter.UserId,
			WorkspaceID: req.Filter.WorkspaceId,
			Statuses: lo.Map(req.Filter.Statuses, func(s pb.BookingStatus, _ int) booking.Status {
				return bookingStatusToModel[s]
			}),
		}

		if req.Filter.DateFrom != nil {
			t := req.Filter.DateFrom.AsTime()
			filter.DateFrom = &t
		}

		if req.Filter.DateTo != nil {
			t := req.Filter.DateTo.AsTime()
			filter.DateTo = &t
		}
	}

	return dto.ListBookingsRequest{
		Page: &page.Page{
			Limit:  req.GetPage().GetLimit(),
			Offset: req.GetPage().GetOffset(),
		},
		Filter: filter,
	}
}

func toListBookingsResponse(resp dto.ListBookingsResponse) *pb.ListBookingsResponse {
	return &pb.ListBookingsResponse{
		Bookings: lo.Map(resp.Bookings, func(b *booking.Booking, _ int) *pb.Booking {
			return toBooking(b)
		}),
		TotalCount: resp.TotalCount,
	}
}

func toBooking(b *booking.Booking) *pb.Booking {
	return &pb.Booking{
		Id:          b.ID,
		WorkspaceId: b.WorkspaceID,
		UserId:      b.UserID,
		StartAt:     timestamppb.New(b.StartAt),
		EndAt:       timestamppb.New(b.EndAt),
		Status:      bookingStatusToProto[b.Status],
		CreatedAt:   timestamppb.New(b.CreatedAt),
		UpdatedAt:   timestamppb.New(b.UpdatedAt),
	}
}
