package booking

import (
	"context"
	"fmt"
	"strconv"

	model "github.com/golangmonster/workspace-booking-service/internal/model/booking"
	"github.com/golangmonster/workspace-booking-service/internal/model/outbox"
)

type workspaceBookingMarshaller func(booking *model.Booking) ([]byte, error)

type service struct {
	repo                  bookingRepository
	userRepo              userRepository
	outboxRepo            outboxRepository
	marshaller            workspaceBookingMarshaller
	workspaceBookingTopic string
}

func New(
	repo bookingRepository,
	userRepo userRepository,
	outboxRepo outboxRepository,
	marshaller workspaceBookingMarshaller,
	workspaceBookingTopic string,
) *service {
	return &service{
		repo:                  repo,
		userRepo:              userRepo,
		outboxRepo:            outboxRepo,
		marshaller:            marshaller,
		workspaceBookingTopic: workspaceBookingTopic,
	}
}

func (s *service) insertWorkspaceBookingOutbox(ctx context.Context, booking *model.Booking) error {
	msg, err := s.marshaller(booking)
	if err != nil {
		return err
	}

	err = s.outboxRepo.CreateOutboxItem(ctx, outbox.Message{
		Topic:       s.workspaceBookingTopic,
		AggregateID: strconv.FormatInt(booking.WorkspaceID, 10),
		Value:       string(msg),
	})
	if err != nil {
		return fmt.Errorf("create outbox: %w", err)
	}

	return nil
}
