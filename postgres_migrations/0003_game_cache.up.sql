
CREATE TABLE game_cache (
  id TEXT PRIMARY KEY,
  title citext,
  series citext,
  developer citext,
  publisher citext,
  release_date TEXT,
  play_mode TEXT[],
  "language" TEXT[],
  extreme BOOLEAN NOT NULL,
  filter_groups TEXT[] NOT NULL,
  original_description citext,
  platform_name citext,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tag_cache (
  id SERIAL PRIMARY KEY,
  "name" citext NOT NULL,
  description TEXT,
  category TEXT NOT NULL,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE game_tag_cache (
  id SERIAL PRIMARY KEY,
  game_id citext NOT NULL,
  tag_id SERIAL NOT NULL,
  CONSTRAINT game_tag_cache_game_id_fkey FOREIGN KEY (game_id) REFERENCES game_cache(id),
  CONSTRAINT game_tag_cache_tag_id_fkey FOREIGN KEY (tag_id) REFERENCES tag_cache(id),
  CONSTRAINT game_tag_cache_game_id_tag_id_key UNIQUE (game_id, tag_id)
);
