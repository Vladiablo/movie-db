CREATE EXTENSION IF NOT EXISTS intarray;

CREATE MATERIALIZED VIEW movies_search_summary AS
SELECT
    id, type, status, title, description, short_description, tagline, age_rating,
    runtime, year_from, release_date, latest_air_date, kp_id, imdb_id,
    genre_ids,
    kp_votes, kp_value, imdb_votes, imdb_value
FROM movies m
CROSS JOIN LATERAL (
    SELECT array_agg(genre_id) AS genre_ids
    FROM movies_genres
    WHERE movie_id = m.id
) g
LEFT JOIN movies_ratings mr ON m.id = mr.movie_id;

CREATE INDEX movies_search_summary_year_from_idx
ON movies_search_summary(year_from);

CREATE INDEX movies_search_summary_type_idx
ON movies_search_summary(type);

CREATE INDEX movies_search_summary_imdb_value_idx
ON movies_search_summary(imdb_value);

CREATE INDEX movies_search_summary_genre_ids_gin_idx
ON movies_search_summary
USING gin(genre_ids gin__int_ops);

CREATE INDEX movies_search_summary_title_trgm_idx
ON movies_search_summary
USING gist(title gist_trgm_ops(siglen=256));

DROP INDEX movies_title_trgm_idx;
