package repository

import (
	"context"
	"movie-db/internal"
)

type Repository interface {
	GetMovieById(ctx context.Context, id int64) (*internal.Movie, error)
	GetSuggestions(ctx context.Context, query string) ([]*internal.MovieSuggestion, error)

	movieRepo()
}
