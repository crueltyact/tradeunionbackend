package news

import "profkom/pkg/s3"

type Service struct {
	repository repository
	s3         *s3.Client
}

func New(repository repository, s3 *s3.Client) *Service {
	return &Service{
		repository: repository,
		s3:         s3,
	}
}
