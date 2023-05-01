CREATE TABLE "transactions" (
  "id" uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  "seller_id" uuid NOT NULL,
  "buyer_id" uuid NOT NULL,
  "product_id" uuid NOT NULL,
  "price_offer" BIGINT NOT NULL,
  "accepted" BOOLEAN DEFAULT FALSE,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted" BOOLEAN DEFAULT FALSE
);

ALTER TABLE "transactions" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");
ALTER TABLE "transactions" ADD FOREIGN KEY ("seller_id") REFERENCES "profiles" ("id");
ALTER TABLE "transactions" ADD FOREIGN KEY ("buyer_id") REFERENCES "profiles" ("id");

CREATE INDEX ON "transactions" ("product_id");
CREATE INDEX ON "transactions" ("buyer_id");
CREATE INDEX ON "transactions" ("seller_id");