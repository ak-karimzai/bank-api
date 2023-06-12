-- +goose Up
CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "fullname" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_last_changed" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

ALTER TABLE "accounts" ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner", "currency");

-- +goose Down
ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "owner_currency_key";
ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "accounts_owner_fkey";
DROP TABLE IF EXISTS "users";