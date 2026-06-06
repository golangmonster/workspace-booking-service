package user

type CreateUserRequest struct {
	Login    string
	FullName string
	Phone    string
}

type UpdateUserRequest struct {
	ID       int64
	Login    string
	FullName string
	Phone    string
}
