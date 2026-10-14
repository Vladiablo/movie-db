package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"movie-db/internal"
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

func (r *PgRepository) GetMovieById(ctx context.Context, id int64) (*internal.Movie, error) {
	const sql = `
SELECT m.id, m.type, m.status, m.title, m.description, m.short_description, m.tagline,
	m.age_rating, m.runtime, m.year_from, m.release_date, m.latest_air_date, g.genre_ids,
	m.kp_id, m.imdb_id, r.ratings, c.crew
FROM
	movies m
	CROSS JOIN LATERAL (
		SELECT array_agg(genre_id) AS genre_ids
		FROM movies_genres
		WHERE movie_id = m.id
	) g
	LEFT JOIN LATERAL (
		SELECT jsonb_build_object(
			'kp', jsonb_build_object(
				'votes', kp_votes,
				'value', kp_value
			),
			'imdb', jsonb_build_object(
				'votes', imdb_votes,
				'value', imdb_value
			)
		) AS ratings
		FROM movies_ratings
		WHERE movie_id = m.id
	) r ON TRUE
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
WHERE id = $1::int
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
