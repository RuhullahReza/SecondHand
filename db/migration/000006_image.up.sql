CREATE TABLE "images" (
  "id" uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  "product_id" uuid NOT NULL,
  "image_url" VARCHAR NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "images" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

CREATE INDEX ON "images" ("product_id");