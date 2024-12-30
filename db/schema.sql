DROP TABLE IF EXISTS movies;
CREATE TABLE movies(
  id SERIAL PRIMARY KEY,
  title VARCHAR(100) NOT NULL,
  year INTEGER NOT NULL,
  runtime INTEGER NOT NULL,
  genres TEXT[] NOT NULL,
  created TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
  version INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX IF NOT EXISTS idx_movies_title ON movies USING GIN (to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS idx_movies_genres ON movies USING GIN (genres);

INSERT INTO movies(title, year, runtime, genres) VALUES('Harry potter', 2005, 220, '{Thriller, Fantasy}'), ('Harry potter 2', 2007, 210, '{Thriller, Fantasy}'), ('Harry potter 2', 2007, 180, '{Drama, Fantasy}'), ('Harry potter 4', 2011, 250, '{Comedy, Thriller, Fantasy}');

DROP TABLE IF EXISTS users;
CREATE TABLE users(
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL,
  password BYTEA NOT NULL,
  active bool NOT NULL,
  created TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
  version INTEGER NOT NULL DEFAULT 1
);
