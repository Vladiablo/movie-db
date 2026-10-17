package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"movie-db/internal"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgRepository struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

var _ Repository = (*PgRepository)(nil)

func NewPgRepository(pool *pgxpool.Pool, logger *slog.Logger) *PgRepository {
	return &PgRepository{pool: pool, logger: logger}
}

const dbTimeout = 10 * time.Second

func (r *PgRepository) movieRepo() {}

// TODO: add mat. view for top crew
func (r *PgRepository) GetMovieById(ctx context.Context, id int64) (*internal.Movie, error) {
	const sql = `
SELECT m.id, m.type, m.status, m.title, m.description, m.short_description, m.tagline,
	m.age_rating, m.runtime, m.year_from, m.release_date, m.latest_air_date,
	m.kp_id, m.imdb_id,
	g.genre_ids,
	mr.kp_votes, mr.kp_value, mr.imdb_votes, mr.imdb_value,
	c.crew
FROM
	movies m
	CROSS JOIN LATERAL (
		SELECT array_agg(genre_id) AS genre_ids
		FROM movies_genres
		WHERE movie_id = m.id
	) g
	LEFT JOIN movies_ratings mr ON mr.movie_id = m.id
	CROSS JOIN LATERAL (
		SELECT jsonb_object_agg(r,
			jsonb_build_object(
				'count', mcc.cnt,
				'members', mcpa.objs
			)
		) AS crew
	FROM
		unnest(enum_range(NULL::role_type)) AS r
		CROSS JOIN LATERAL (
			SELECT jsonb_agg(mcp.obj) AS objs
			FROM (
				SELECT jsonb_build_object(
					'id', mc.person_id,
					'name', p.name,
					'role_details', mc.role_details
				) AS obj
				FROM movies_crew mc
				JOIN persons p ON mc.person_id = p.id
				WHERE mc.movie_id = m.id AND mc.role = r
				ORDER BY mc.ordering ASC
				LIMIT 5
			) mcp
		) mcpa
		CROSS JOIN LATERAL (
			SELECT COUNT(*) AS cnt
			FROM movies_crew mc
			WHERE mc.movie_id = m.id AND mc.role = r
		) mcc
		WHERE mcpa.objs IS NOT NULL
	) c
WHERE m.id = $1::int
`

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	rows, err := r.pool.Query(ctx, sql, id)
	if err != nil {
		r.logger.Error("Failed to query movie",
			slog.Int64("id", id),
			slog.Any("error", err),
		)

		return nil, fmt.Errorf("failed to query movie: %w", err)
	}

	movie, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[pgMovie])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		r.logger.Error("Failed to scan movie",
			slog.Int64("id", id),
			slog.Any("error", err),
		)

		return nil, fmt.Errorf("failed to scan movie: %w", err)
	}

	return movie.copyToModel(), nil
}

// TODO: Also order by weighted rating
func (r *PgRepository) GetSuggestions(ctx context.Context, query string) ([]*internal.MovieSuggestion, error) {
	const sql = `
SELECT
    id, type, status, title, description, short_description, tagline, age_rating,
    runtime, year_from, release_date, latest_air_date, kp_id, imdb_id,
    genre_ids,
    imdb_value, imdb_votes, kp_value, kp_votes 
FROM movies_search_summary
WHERE title % $1::text
ORDER BY title <-> $1::text ASC
LIMIT 10;
`

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	rows, err := r.pool.Query(ctx, sql, query)
	if err != nil {
		r.logger.Error("Failed to query suggestions",
			slog.String("query", query),
			slog.Any("error", err),
		)

		return nil, fmt.Errorf("failed to query suggestions: %w", err)
	}

	suggestions, err := pgx.CollectRows(rows, pgx.RowToStructByName[pgMovieSuggestion])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		r.logger.Error("Failed to scan movie suggestions",
			slog.String("query", query),
			slog.Any("error", err),
		)

		return nil, fmt.Errorf("failed to scan movie suggestions: %w", err)
	}

	suggestionsCopy := make([]*internal.MovieSuggestion, len(suggestions))
	for i := range suggestionsCopy {
		suggestionsCopy[i] = suggestions[i].copyToModel()
	}

	return suggestionsCopy, nil
}

func (r *PgRepository) SearchMovies(ctx context.Context, params *MoviesSearchParams) ([]*internal.MovieSuggestion, error) {
	const sql = `
SELECT
    id, type, status, title, description, short_description, tagline, age_rating,
    runtime, year_from, release_date, latest_air_date, kp_id, imdb_id,
    genre_ids,
    imdb_value, imdb_votes, kp_value, kp_votes 
FROM movies_search_summary
WHERE
    ($1::movie_type IS NULL OR type = $1::movie_type) 
    AND ($2::int[] IS NULL OR genre_ids @> $2::int[])
    AND ($3::float4 = 0.0 OR imdb_value >= $3::float4)
    AND ($4::float4 = 0.0 OR imdb_value <= $4::float4)
    AND ($5::int = 0 OR year_from >= $5::int)
    AND ($6::int = 0 OR year_from <= $6::int)
`

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

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

	queryBuilder.WriteString("LIMIT 50;")

	rows, err := r.pool.Query(ctx, queryBuilder.String(),
		params.Type, params.GenreIds,
		params.ImdbRatingFrom, params.ImdbRatingTo,
		params.YearFrom, params.YearTo,
	)
	if err != nil {
		r.logger.Error("Failed to search movies",
			slog.Any("type", params.Type),
			slog.Any("genreIds", params.GenreIds),
			slog.Float64("imdbRatingFrom", float64(params.ImdbRatingFrom)),
			slog.Float64("imdbRatingTo", float64(params.ImdbRatingTo)),
			slog.Int64("yearFrom", int64(params.YearFrom)),
			slog.Int64("yearTo", int64(params.YearTo)),
			slog.String("orderBy", params.OrderBy),
			slog.String("orderDir", params.OrderDir),
			slog.Any("error", err),
		)

		return nil, fmt.Errorf("failed to search movies: %w", err)
	}

	suggestions, err := pgx.CollectRows(rows, pgx.RowToStructByName[pgMovieSuggestion])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		r.logger.Error("Failed to scan movies search results",
			slog.Any("type", params.Type),
			slog.Any("genreIds", params.GenreIds),
			slog.Float64("imdbRatingFrom", float64(params.ImdbRatingFrom)),
			slog.Float64("imdbRatingTo", float64(params.ImdbRatingTo)),
			slog.Int64("yearFrom", int64(params.YearFrom)),
			slog.Int64("yearTo", int64(params.YearTo)),
			slog.String("orderBy", params.OrderBy),
			slog.String("orderDir", params.OrderDir),
			slog.Any("error", err),
		)

		return nil, fmt.Errorf("failed to scan movies search results: %w", err)
	}

	suggestionsCopy := make([]*internal.MovieSuggestion, len(suggestions))
	for i := range suggestionsCopy {
		suggestionsCopy[i] = suggestions[i].copyToModel()
	}

	return suggestionsCopy, nil
}
