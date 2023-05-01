CREATE TABLE "cities" (
  "id" uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  "name" VARCHAR NOT NULL UNIQUE,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);