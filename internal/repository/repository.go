package repository

import (
	"context"
	"movie-db/internal"
)

type Repository interface {
	GetMovieById(ctx context.Context, id int64) (*internal.Movie, error)
	GetSuggestions(ctx context.Context, query string) ([]*internal.MovieSuggestion, error)
	SearchMovies(ctx context.Context, params *MoviesSearchParams) ([]*internal.MovieSuggestion, error)

	movieRepo()
}

type MoviesSearchParams struct {
	Type           *string
	YearFrom       int32
	YearTo         int32
	GenreIds       []int32
	ImdbRatingFrom float32
	ImdbRatingTo   float32
	OrderBy        string
	OrderDir       string
}
