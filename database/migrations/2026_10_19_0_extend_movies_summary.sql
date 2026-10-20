DROP MATERIALIZED VIEW movies_search_summary;

CREATE MATERIALIZED VIEW movies_summary AS
SELECT
    id, type, status, title, description, short_description, tagline, age_rating,
    runtime, year_from, release_date, latest_air_date, kp_id, imdb_id,
    genre_ids,
    kp_votes, kp_value, imdb_votes, imdb_value,
    tc.crew AS top_crew, c.crew AS crew
FROM
    movies m
    CROSS JOIN LATERAL (
        SELECT array_agg(genre_id) AS genre_ids
        FROM movies_genres
        WHERE movie_id = m.id
    ) g
    LEFT JOIN movies_ratings mr ON m.id = mr.movie_id
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
    ) tc
    CROSS JOIN LATERAL (
        SELECT array_agg(person_id) AS crew
        FROM movies_crew
        WHERE movie_id = m.id
    ) c;

CREATE INDEX movies_summary_id_idx
ON movies_summary(id);

CREATE INDEX movies_summary_year_from_idx
ON movies_summary(year_from);

CREATE INDEX movies_summary_type_idx
ON movies_summary(type);

CREATE INDEX movies_summary_imdb_value_idx
ON movies_summary(imdb_value);

CREATE INDEX movies_summary_genre_ids_gin_idx
ON movies_summary
USING gin(genre_ids gin__int_ops);

CREATE INDEX movies_summary_title_trgm_idx
ON movies_summary
USING gist(title gist_trgm_ops(siglen=256));

CREATE INDEX movies_summary_crew_idx
ON movies_summary
USING gin(crew gin__int_ops);
