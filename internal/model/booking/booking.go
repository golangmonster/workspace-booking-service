package booking

import (
	"errors"
	"time"
)

var (
	ErrBookingNotFound  = errors.New("booking not found")
	ErrBookingOverlap   = errors.New("workspace is already booked for this time slot")
	ErrInvalidTimeRange = errors.New("start time must be before end time")
)

type Status string

const (
	StatusUnspecified Status = "UNSPECIFIED"
	StatusActive      Status = "ACTIVE"
	StatusCancelled   Status = "CANCELLED"
	StatusCompleted   Status = "COMPLETED"
)

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
