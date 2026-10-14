ALTER TABLE imdb.titles_ratings
ALTER COLUMN avg_rating TYPE float4;

ALTER TABLE movies_ratings
ALTER COLUMN kp_votes TYPE bigint;

ALTER TABLE movies_ratings
ALTER COLUMN kp_value TYPE float4;

ALTER TABLE movies_ratings
ALTER COLUMN imdb_votes TYPE bigint;

ALTER TABLE movies_ratings
ALTER COLUMN imdb_value TYPE float4;
