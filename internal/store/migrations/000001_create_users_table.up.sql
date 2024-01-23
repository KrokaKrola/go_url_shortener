CREATE TABLE IF NOT EXISTS links(
  id serial PRIMARY KEY,
  url text NOT NULL,
  hash text NOT NULL,
  created_at timestamp NOT NULL DEFAULT NOW(),
  rate_limit_per_minute integer NOT NULL DEFAULT 60
);