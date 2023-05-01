CREATE TABLE "profiles" (
  "id" uuid PRIMARY KEY,
  "name" VARCHAR NOT NULL,
  "city" VARCHAR,
  "address" VARCHAR,
  "phone_number" VARCHAR,
  "image_url" VARCHAR,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "profiles" ADD FOREIGN KEY ("id") REFERENCES "accounts" ("id");