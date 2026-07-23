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
