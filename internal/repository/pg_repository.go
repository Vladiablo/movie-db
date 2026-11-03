package repository

import (
	"context"
	"errors"
	"fmt"
	"movie-db/internal"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgRepository struct {
	pool      *pgxpool.Pool
	dbTimeout time.Duration
}

var _ Repository = (*PgRepository)(nil)

func NewPgRepository(pool *pgxpool.Pool, dbTimeout time.Duration) *PgRepository {
	return &PgRepository{
		pool:      pool,
		dbTimeout: max(0, dbTimeout),
	}
}

func (r *PgRepository) movieRepo() {}

func (r *PgRepository) GetMovieById(ctx context.Context, id int64) (*internal.Movie, error) {
	const sql = `
SELECT
    id, type, status, title, description, short_description, tagline,
    age_rating, runtime, year_from, release_date, latest_air_date,
    kp_id, imdb_id,
    genre_ids,
    kp_votes, kp_value, imdb_votes, imdb_value,
    top_crew AS crew
FROM movies_summary 
WHERE id = $1::int
`

	if r.dbTimeout > 0 {
		c, cancel := context.WithTimeout(ctx, r.dbTimeout)
		defer cancel()

		ctx = c
	}

	rows, err := r.pool.Query(ctx, sql, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query movie: %w", err)
	}

	movie, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[pgMovie])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to scan movie: %w", err)
	}

	return movie.toModel(), nil
}

// TODO: Also order by weighted rating
func (r *PgRepository) GetSuggestions(ctx context.Context, query string) ([]*internal.MovieBase, error) {
	const sql = `
SELECT
    id, type, status, title, description, short_description, tagline, age_rating,
    runtime, year_from, release_date, latest_air_date, kp_id, imdb_id,
    genre_ids,
    imdb_value, imdb_votes, kp_value, kp_votes 
FROM movies_summary
WHERE title % $1::text
ORDER BY title <-> $1::text ASC
LIMIT 10;
`

	if r.dbTimeout > 0 {
		c, cancel := context.WithTimeout(ctx, r.dbTimeout)
		defer cancel()

		ctx = c
	}

	rows, err := r.pool.Query(ctx, sql, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query suggestions: %w", err)
	}

	suggestions, err := pgx.CollectRows(rows, pgx.RowToStructByName[pgMovieBase])
	if err != nil {
		return nil, fmt.Errorf("failed to scan movie suggestions: %w", err)
	}

	suggestionsCopy := make([]*internal.MovieBase, len(suggestions))
	for i := range suggestionsCopy {
		suggestionsCopy[i] = suggestions[i].toModel()
	}

	return suggestionsCopy, nil
}

func (r *PgRepository) SearchMovies(ctx context.Context, params *MoviesSearchParams) ([]*internal.MovieBase, error) {
	const sql = `
SELECT
    id, type, status, title, description, short_description, tagline, age_rating,
    runtime, year_from, release_date, latest_air_date, kp_id, imdb_id,
    genre_ids,
    imdb_value, imdb_votes, kp_value, kp_votes 
FROM movies_summary
WHERE
    ($1::movie_type IS NULL OR type = $1::movie_type) 
    AND ($2::int[] IS NULL OR genre_ids @> $2::int[])
    AND ($3::float4 = 0.0 OR imdb_value >= $3::float4)
    AND ($4::float4 = 0.0 OR imdb_value <= $4::float4)
    AND ($5::int = 0 OR year_from >= $5::int)
    AND ($6::int = 0 OR year_from <= $6::int)
    AND ($7::int = 0 OR crew @@ $7::text::query_int)
`

	if r.dbTimeout > 0 {
		c, cancel := context.WithTimeout(ctx, r.dbTimeout)
		defer cancel()

		ctx = c
	}

	orderBy := params.OrderBy
	if orderBy == "imdb_rating" {
		orderBy = "imdb_value"
	}

	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(sql)

	queryBuilder.WriteString("ORDER BY ")
	queryBuilder.WriteString(orderBy)

	queryBuilder.WriteByte(' ')
	queryBuilder.WriteString(params.OrderDir)
	queryBuilder.WriteByte('\n')

	queryBuilder.WriteString("LIMIT 50")

	rows, err := r.pool.Query(ctx, queryBuilder.String(),
		params.Type, params.GenreIds,
		params.ImdbRatingFrom, params.ImdbRatingTo,
		params.YearFrom, params.YearTo,
		params.PersonId,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search movies: %w", err)
	}

	searchResults, err := pgx.CollectRows(rows, pgx.RowToStructByName[pgMovieBase])
	if err != nil {
		return nil, fmt.Errorf("failed to scan movies search results: %w", err)
	}

	result := make([]*internal.MovieBase, len(searchResults))
	for i := range result {
		result[i] = searchResults[i].toModel()
	}

	return result, nil
}
