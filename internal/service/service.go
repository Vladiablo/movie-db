package service

import (
	"context"
	"movie-db/internal"
	"movie-db/internal/repository"
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetMovieById(ctx context.Context, id int64) (*internal.Movie, error) {
	return s.repo.GetMovieById(ctx, id)
}
