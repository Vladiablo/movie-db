CREATE SCHEMA IF NOT EXISTS imdb;

CREATE TABLE IF NOT EXISTS imdb.titles_basic (
    tconst text PRIMARY KEY,
    type text NOT NULL,
    primary_title text NOT NULL,
    original_title text NOT NULL,
    is_adult bool NOT NULL,
    start_year int,
    end_year int,
    runtime_minutes int,
    genres text
);

CREATE TABLE IF NOT EXISTS imdb.titles_ratings (
    tconst text PRIMARY KEY,
    avg_rating float NOT NULL,
    num_votes bigint NOT NULL
);

CREATE TABLE IF NOT EXISTS imdb.name_basics (
    nconst text PRIMARY KEY,
    primary_name text,
    birth_year int,
    death_year int,
    primary_profession text,
    known_for_titles text
);

CREATE TABLE IF NOT EXISTS imdb.title_crew (
    tconst text PRIMARY KEY,
    directors text,
    writers text
);

CREATE TABLE IF NOT EXISTS imdb.title_principals (
    tconst text,
    ordering int,
    nconst text NOT NULL,
    category text NOT NULL,
    job text,
    characters text[],

    PRIMARY KEY (tconst, ordering)
);

CREATE TABLE IF NOT EXISTS imdb.persons (
    id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name text,
    nconst text NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS imdb.movies (
    id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    type text,
    title text NOT NULL,
    runtime int CHECK (runtime >= 0),
    year_from int CHECK (year_from > 0),
    tconst text NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS imdb.genres (
    id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    slug text NOT NULL UNIQUE,
    name text NOT NULL
);

CREATE TABLE IF NOT EXISTS imdb.movies_genres (
    movie_id int NOT NULL REFERENCES imdb.movies(id) ON DELETE CASCADE,
    genre_id int NOT NULL REFERENCES imdb.genres(id) ON DELETE CASCADE,

    PRIMARY KEY (movie_id, genre_id)
);

CREATE TABLE IF NOT EXISTS imdb.movies_crew (
    id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    movie_id int NOT NULL REFERENCES imdb.movies(id) ON DELETE CASCADE,
    person_id int NOT NULL REFERENCES imdb.persons(id) ON DELETE CASCADE,
    ordering int NOT NULL CHECK (ordering > 0),
    role role_type NOT NULL, -- actor, composer, etc
    role_details text,

    UNIQUE (movie_id, person_id, role)
);
