package guide

type Service struct {
	repo repository
}

func New(repo repository) *Service {
	return &Service{
		repo: repo,
	}
}
