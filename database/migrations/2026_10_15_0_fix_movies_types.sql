CREATE TABLE IF NOT EXISTS tv_episodes (
    id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    status text NOT NULL DEFAULT 'post_production',
    title text NOT NULL,
    description text,
    short_description text,
    tagline text,
    age_rating text,
    runtime int CHECK (runtime >= 0),
    year_from int CHECK (year_from > 0),
    release_date date,
    latest_air_date date,
    kp_id int UNIQUE CHECK (kp_id > 0),
    imdb_id int UNIQUE CHECK (imdb_id > 0)
);

CREATE TABLE IF NOT EXISTS tv_episodes_ratings (
    tv_episode_id int PRIMARY KEY REFERENCES movies(id) ON DELETE CASCADE,
    kp_votes bigint,
    kp_value float4,
    imdb_votes bigint,
    imdb_value float4
);

WITH inserted_episode AS (
    INSERT INTO tv_episodes (status, title, description, short_description, tagline, age_rating, runtime, year_from, release_date, latest_air_date, kp_id, imdb_id)
        SELECT status, title, description, short_description, tagline, age_rating, runtime, year_from, release_date, latest_air_date, kp_id, imdb_id
        FROM movies m
        WHERE type = 'tv_episode'
    RETURNING id AS new_id, imdb_id
)
INSERT INTO tv_episodes_ratings (tv_episode_id, kp_votes, kp_value, imdb_votes, imdb_value)
SELECT e.new_id, mr.kp_votes, mr.kp_value, mr.imdb_votes, mr.imdb_value
FROM inserted_episode e
JOIN movies m ON e.imdb_id = m.imdb_id
JOIN movies_ratings mr ON m.id = mr.movie_id;

DELETE FROM movies
WHERE type = 'tv_episode';

DROP INDEX movies_title_trgm_idx;

UPDATE movies
SET type = 'movie'
WHERE type NOT IN ('tv_mini_series', 'tv_series');

ALTER TABLE movies
ALTER COLUMN type DROP DEFAULT,
ALTER COLUMN type TYPE text;

DROP TYPE movie_type;

CREATE TYPE movie_type AS ENUM (
    'movie', 'series'
);

ALTER TABLE movies
ALTER COLUMN type TYPE movie_type USING (
    CASE type
        WHEN 'movie' THEN 'movie'::movie_type
        ELSE 'series'::movie_type
    END
),
ALTER COLUMN type SET DEFAULT 'movie';

CREATE INDEX movies_title_trgm_idx ON movies
USING gist(title gist_trgm_ops(siglen=256));
