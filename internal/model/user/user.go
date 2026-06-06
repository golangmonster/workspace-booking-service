package user

import (
	"errors"
	"time"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrLoginExists  = errors.New("user with this login already exists")
	ErrPhoneExists  = errors.New("user with this phone already exists")
)

type User struct {
	ID        int64
	Login     string
	FullName  string
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
	IsDeleted bool
}
