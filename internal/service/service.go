package service

import (
	"context"
	"movie-db/internal"
	"movie-db/internal/helpers"
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

type MoviesSearchParams struct {
	Type           string
	YearFrom       int32
	YearTo         int32
	GenreIds       []int32
	ImdbRatingFrom float32
	ImdbRatingTo   float32
	OrderBy        string
	OrderDir       string
}

func (s *Service) GetMovieById(ctx context.Context, id int64) (*internal.Movie, error) {
	return s.repo.GetMovieById(ctx, id)
}

func (s *Service) GetSuggestions(ctx context.Context, query string) ([]*internal.MovieSuggestion, error) {
	return s.repo.GetSuggestions(ctx, query)
}

func (s *Service) SearchMovies(ctx context.Context, params *MoviesSearchParams) ([]*internal.MovieSuggestion, error) {
	movieType := (*string)(nil)
	if params.Type != "" {
		movieType = &params.Type
	}

	orderBy := params.OrderBy
	if orderBy == "" {
		orderBy = "imdb_rating"
	}

	orderDir := params.OrderDir
	if orderDir == "" {
		orderDir = "desc"
	}

	return s.repo.SearchMovies(ctx, &repository.MoviesSearchParams{
		Type:           movieType,
		YearFrom:       params.YearFrom,
		YearTo:         params.YearTo,
		GenreIds:       helpers.CopySlice(params.GenreIds),
		ImdbRatingFrom: params.ImdbRatingFrom,
		ImdbRatingTo:   params.ImdbRatingTo,
		OrderBy:        orderBy,
		OrderDir:       orderDir,
	})
}
