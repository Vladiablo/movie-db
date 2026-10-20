package service

import (
	"context"
	"log/slog"
	"movie-db/internal"
	"movie-db/internal/repository"
	"slices"
)

type Service struct {
	repo repository.Repository

	logger *slog.Logger
}

func New(repo repository.Repository, logger *slog.Logger) *Service {
	return &Service{
		repo: repo, logger: logger,
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
	PersonId       int32
}

func (s *Service) GetMovieById(ctx context.Context, id int64) (*internal.Movie, error) {
	result, err := s.repo.GetMovieById(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get movie",
			slog.Any("error", err),
			slog.Int64("id", id),
		)
	}

	return result, err
}

func (s *Service) GetSuggestions(ctx context.Context, query string) ([]*internal.MovieBase, error) {
	result, err := s.repo.GetSuggestions(ctx, query)
	if err != nil {
		s.logger.Error("Failed to get suggestions",
			slog.Any("error", err),
			slog.String("query", query),
		)
	}

	return result, err
}

func (s *Service) SearchMovies(ctx context.Context, params *MoviesSearchParams) ([]*internal.MovieBase, error) {
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

	result, err := s.repo.SearchMovies(ctx, &repository.MoviesSearchParams{
		Type:           movieType,
		YearFrom:       params.YearFrom,
		YearTo:         params.YearTo,
		GenreIds:       slices.Clone(params.GenreIds),
		ImdbRatingFrom: params.ImdbRatingFrom,
		ImdbRatingTo:   params.ImdbRatingTo,
		OrderBy:        orderBy,
		OrderDir:       orderDir,
		PersonId:       params.PersonId,
	})

	if err != nil {
		s.logger.Error("Failed to search movies",
			slog.Any("error", err),
			slog.Any("type", movieType),
			slog.Any("genreIds", params.GenreIds),
			slog.Float64("imdbRatingFrom", float64(params.ImdbRatingFrom)),
			slog.Float64("imdbRatingTo", float64(params.ImdbRatingTo)),
			slog.Int64("yearFrom", int64(params.YearFrom)),
			slog.Int64("yearTo", int64(params.YearTo)),
			slog.String("orderBy", orderBy),
			slog.String("orderDir", orderDir),
			slog.Int64("personId", int64(params.PersonId)),
		)
	}

	return result, err
}
