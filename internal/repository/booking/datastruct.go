package booking

import "time"

type bookingItem struct {
	ID          int64     `db:"id"`
	WorkspaceID int64     `db:"workspace_id"`
	UserID      int64     `db:"user_id"`
	StartAt     time.Time `db:"start_at"`
	EndAt       time.Time `db:"end_at"`
	Status      string    `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type workspaceItem struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	Lat         float64   `db:"lat"`
	Lon         float64   `db:"lon"`
	FullAddress string    `db:"full_address"`
	Type        string    `db:"type"`
	Status      string    `db:"status"`
	Capacity    uint32    `db:"capacity"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	IsDeleted   bool      `db:"is_deleted"`
}
