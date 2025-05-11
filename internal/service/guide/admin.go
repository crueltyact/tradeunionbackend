package guide

import "context"

func (s *Service) DeleteGuide(ctx context.Context, id int) (err error) {
	return s.repo.DeleteGuide(ctx, id)
}

func (s *Service) UpdateGuide(ctx context.Context) (err error) {
	return err
}

func (s *Service) DeleteTheme(ctx context.Context, id int) (err error) {
	return s.repo.DeleteTheme(ctx, id)
}
