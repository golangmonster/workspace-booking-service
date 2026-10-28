package booking

import (
	"errors"
	"slices"
	"time"
)

var (
	ErrBookingNotFound         = errors.New("booking not found")
	ErrBookingOverlap          = errors.New("workspace is already booked for this time slot")
	ErrInvalidTimeRange        = errors.New("start time must be before end time")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
)

type Status string

const (
	StatusUnspecified Status = "UNSPECIFIED"
	StatusActive      Status = "ACTIVE"
	StatusCancelled   Status = "CANCELLED"
	StatusCompleted   Status = "COMPLETED"
)

var (
	allowedStatusTransitions = map[Status][]Status{
		StatusActive:    {StatusCancelled, StatusCompleted},
		StatusCancelled: {},
		StatusCompleted: {},
	}
)

func StatusTransitionAllowed(from, to Status) bool {
	return slices.Contains(allowedStatusTransitions[from], to)
}

type Booking struct {
	ID          int64
	WorkspaceID int64
	UserID      int64
	StartAt     time.Time
	EndAt       time.Time
	Status      Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
