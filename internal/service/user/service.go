package user

type service struct {
	repo userRepository
}

func New(repo userRepository) *service {
	return &service{
		repo: repo,
	}
}
