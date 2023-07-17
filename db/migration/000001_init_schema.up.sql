CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "balance" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "country_code" int
);

CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint not null,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  FOREIGN KEY ("account_id") REFERENCES "accounts" ("id")
);

CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "from_account_id" bigint not null,
  "to_account_id" bigint not null,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id"),
  FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id")
);

CREATE INDEX ON "accounts" ("owner");
CREATE INDEX ON "entries" ("account_id");
CREATE INDEX ON "transfers" ("from_account_id");
CREATE INDEX ON "transfers" ("to_account_id");
CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");
