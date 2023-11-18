CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE "fpcomm_role" (
  "id" TEXT PRIMARY KEY,
  "name" citext NOT NULL,
  "color" TEXT NOT NULL
);

CREATE TABLE "fpcomm_user" (
  "id" TEXT PRIMARY KEY,
  "name" citext NOT NULL,
  "avatar" TEXT,
  "roles" TEXT[],
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "session" (
  "id" SERIAL PRIMARY KEY,
  "uid" TEXT NOT NULL,
  "secret" TEXT NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "expires_at" TIMESTAMP NOT NULL,
  "ip_addr" TEXT NOT NULL
);
