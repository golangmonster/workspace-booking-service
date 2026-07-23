package workspace

import (
	"errors"
	"time"
)

type (
	Type   string
	Status string
)

const (
	TypeUnspecified Type = "UNSPECIFIED"
	TypeDesk        Type = "DESK"
	TypeMeetingRoom Type = "MEETING_ROOM"
	TypePhoneBooth  Type = "PHONE_BOOTH"

	StatusUnspecified Status = "UNSPECIFIED"
	StatusAvailable   Status = "AVAILABLE"
	StatusBooked      Status = "BOOKED"
	StatusMaintenance Status = "MAINTENANCE"
)

var (
	ErrWorkspaceNotFound      = errors.New("workspace not found")
	ErrWorkspaceAlreadyExists = errors.New("workspace with such name already exists in this type")
	ErrWorkspaceNotAvailable  = errors.New("workspace is not available for booking")
)

type Workspace struct {
	ID          int64
	Name        string
	Lat         float64
	Lon         float64
	FullAddress string
	Type        Type
	Status      Status
	Capacity    uint32
	CreatedAt   time.Time
	UpdatedAt   time.Time
	IsDeleted   bool
}
