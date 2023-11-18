CREATE TABLE playlist (
  id SERIAL PRIMARY KEY,
  name citext NOT NULL,
  total_games INTEGER NOT NULL DEFAULT 0,
  description citext NOT NULL,
  author_id TEXT NOT NULL REFERENCES fpcomm_user(id),
  icon TEXT,
  library TEXT NOT NULL,
  public BOOLEAN NOT NULL DEFAULT FALSE,
  extreme BOOLEAN NOT NULL,
  filter_groups TEXT[] NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX playlists_name_idx ON playlist(name);
CREATE INDEX playlists_total_games_idx ON playlist(total_games);
CREATE INDEX playlists_created_at_idx ON playlist(created_at);
CREATE INDEX playlists_updated_at_idx ON playlist(updated_at);

CREATE TABLE playlist_game (
  playlist_id INTEGER NOT NULL REFERENCES playlist(id),
  game_id TEXT NOT NULL,
  notes citext,
  PRIMARY KEY (playlist_id, game_id)
);

CREATE INDEX playlist_game_playlist_id_idx ON playlist_game(playlist_id);
