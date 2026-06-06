package user

import "github.com/golangmonster/workspace-booking-service/internal/model/user"

func toUser(i userItem) *user.User {
	return &user.User{
		ID:        i.ID,
		Login:     i.Login,
		FullName:  i.FullName,
		Phone:     i.Phone,
		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
		IsDeleted: i.IsDeleted,
	}
}
