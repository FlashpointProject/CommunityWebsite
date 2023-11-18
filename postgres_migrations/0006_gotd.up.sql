CREATE TABLE "gotd_suggestion" (
  "id" SERIAL PRIMARY KEY,
  "game_id" citext NOT NULL,
  "author_id" TEXT NOT NULL,
  "anonymous" BOOLEAN NOT NULL,
  "description" citext NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
  "suggested_date" DATE
);

