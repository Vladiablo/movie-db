-- TODO: Maybe add enums
-- TODO: Maybe add indices, more constraints

CREATE TYPE movie_type AS ENUM (
    'movie', 'short', 'tv_episode', 'tv_mini_series',
    'tv_movie', 'tv_pilot', 'tv_series', 'tv_short',
    'tv_special', 'video'
);

CREATE TYPE role_type AS ENUM (
    'actor', 'archive_footage', 'archive_sound', 'casting_director',
    'cinematographer', 'composer', 'director', 'editor',
    'producer', 'production_designer', 'self', 'writer'
);

CREATE TABLE IF NOT EXISTS movies (
    id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    type movie_type NOT NULL DEFAULT 'movie',
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

CREATE TABLE IF NOT EXISTS movies_posters (
    movie_id int PRIMARY KEY REFERENCES movies(id) ON DELETE CASCADE,
    type text NOT NULL,
    path text NOT NULL,
    width int CHECK (width > 0),
    height int CHECK (height > 0)
);

CREATE TABLE IF NOT EXISTS genres (
    id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    slug text NOT NULL UNIQUE,
    name text NOT NULL,
    is_hidden bool NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS movies_genres (
    movie_id int NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    genre_id int NOT NULL REFERENCES genres(id) ON DELETE CASCADE,

    PRIMARY KEY (movie_id, genre_id)
);

CREATE TABLE IF NOT EXISTS countries (
    id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    iso_3166_1 text NOT NULL,
    name text NOT NULL
);

CREATE TABLE IF NOT EXISTS movies_countries (
    movie_id int NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    country_id int NOT NULL REFERENCES countries(id) ON DELETE CASCADE,

    PRIMARY KEY (movie_id, country_id)
);

CREATE TABLE IF NOT EXISTS balancer_translations (
    id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    balancer text NOT NULL, -- aloha, rezka etc
    balancer_id int NOT NULL,
    name text NOT NULL,

    UNIQUE (balancer, balancer_id)
);

CREATE TABLE IF NOT EXISTS movies_watch_options (
    id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    balancer text NOT NULL, -- aloha, rezka, etc
    balancer_id text NOT NULL,
    balancer_type text NOT NULL,
    movie_id int NOT NULL REFERENCES movies(id) ON DELETE CASCADE,

    UNIQUE (balancer, balancer_id)
);

CREATE TABLE IF NOT EXISTS movies_watch_options_translations (
    movie_watch_option_id int NOT NULL REFERENCES movies_watch_options(id) ON DELETE CASCADE,
    translation_id int NOT NULL REFERENCES balancer_translations(id) ON DELETE CASCADE,

    PRIMARY KEY (movie_watch_option_id, translation_id)
);

-- Maybe split into kp_rating and imdb_ratings
CREATE TABLE IF NOT EXISTS movies_ratings (
    movie_id int PRIMARY KEY REFERENCES movies(id) ON DELETE CASCADE,
    kp_votes int,
    kp_value float,
    imdb_votes int,
    imdb_value float
);

CREATE TABLE IF NOT EXISTS persons (
    id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name text NOT NULL
);

CREATE TABLE IF NOT EXISTS movies_crew (
    id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    movie_id int NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    person_id int NOT NULL REFERENCES persons(id) ON DELETE CASCADE,
    ordering int NOT NULL CHECK (ordering > 0),
    role role_type NOT NULL, -- actor, composer, etc
    role_details text,

    UNIQUE (movie_id, person_id, role)
);

ALTER TABLE movies_crew
ALTER COLUMN role TYPE role_type USING role::role_type;


/*
CREATE TABLE IF NOT EXISTS networks (
    id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    slug text NOT NULL UNIQUE,
    name text NOT NULL
);

CREATE TABLE IF NOT EXISTS movies_networks (
    movie_id int NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    network_id int NOT NULL REFERENCES networks(id) ON DELETE CASCADE,

    PRIMARY KEY (movie_id, network_id)
);

CREATE TABLE IF NOT EXISTS interest_categories (
    id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    slug text NOT NULL UNIQUE,
    title text NOT NULL
);

CREATE TABLE IF NOT EXISTS interest_kinds (
    id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    slug text NOT NULL UNIQUE,
    title text NOT NULL
);

CREATE TABLE IF NOT EXISTS interests (
    id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    slug text NOT NULL UNIQUE,
    category_id int NOT NULL REFERENCES interest_categories(id),
    kind_id int NOT NULL REFERENCES interest_kinds(id),
    title text NOT NULL,
    description text
);

-- Maybe merge with movies_posters
CREATE TABLE IF NOT EXISTS interests_posters (
    interest_id int PRIMARY KEY REFERENCES interests(id) ON DELETE CASCADE,
    type text NOT NULL,
    path text NOT NULL,
    width int CHECK (width > 0),
    height int check (width > 0)
);

CREATE TABLE IF NOT EXISTS movies_interests (
    movie_id int NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    interest_id int NOT NULL REFERENCES interests(id) ON DELETE CASCADE,

    PRIMARY KEY (movie_id, interest_id)
);
 */
