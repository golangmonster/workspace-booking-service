package booking

type service struct {
	repo     bookingRepository
	userRepo userRepository
}

func New(repo bookingRepository, userRepo userRepository) *service {
	return &service{
		repo:     repo,
		userRepo: userRepo,
	}
}
