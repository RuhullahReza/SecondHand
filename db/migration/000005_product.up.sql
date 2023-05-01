CREATE TABLE "products" (
  "id" uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  "account_id" uuid,
  "name" VARCHAR NOT NULL,
  "price" BIGINT NOT NULL,
  "category" VARCHAR NOT NULL,
  "description" VARCHAR NOT NULL,
  "thumbnail" VARCHAR,
  "sold" BOOLEAN NOT NULL DEFAULT FALSE,
  "published" BOOLEAN NOT NULL DEFAULT FALSE,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted" BOOLEAN NOT NULL DEFAULT FALSE
);

ALTER TABLE "products" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

CREATE INDEX ON "products" ("account_id");

CREATE INDEX ON "products" ("category");