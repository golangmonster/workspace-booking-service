package booking

type service struct {
	repo bookingRepository
}

func New(repo bookingRepository) *service {
	return &service{
		repo: repo,
	}
}
