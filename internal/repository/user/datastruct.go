package user

import "time"

type userItem struct {
	ID        int64     `db:"id"`
	Login     string    `db:"login"`
	FullName  string    `db:"full_name"`
	Phone     string    `db:"phone"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	IsDeleted bool      `db:"is_deleted"`
}
