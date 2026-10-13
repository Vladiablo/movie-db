CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX movies_title_trgm_idx ON movies
USING gist(title gist_trgm_ops(siglen=256))
WHERE type != 'tv_episode';
