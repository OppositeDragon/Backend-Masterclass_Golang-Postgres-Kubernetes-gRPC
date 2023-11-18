DROP TABLE IF EXISTS "transfer";
DROP TABLE IF EXISTS "entry";
DROP TABLE IF EXISTS "account";

CREATE TABLE "account" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "currency" varchar NOT NULL,
  "balance" bigint NOT NULL,
  "createdAt" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "entry" (
  "id" bigserial PRIMARY KEY,
  "accountId" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "createdAt" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transfer" (
  "id" bigserial PRIMARY KEY,
  "fromAccountId" bigint NOT NULL,
  "toAccountId" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "createdAt" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "account" ("owner");

CREATE INDEX ON "entry" ("accountId");

CREATE INDEX ON "transfer" ("fromAccountId");

CREATE INDEX ON "transfer" ("toAccountId");

CREATE INDEX ON "transfer" ("fromAccountId", "toAccountId");

CREATE INDEX ON "transfer" ("toAccountId", "createdAt");

COMMENT ON COLUMN "entry"."amount" IS 'can be positive or negative';

COMMENT ON COLUMN "transfer"."amount" IS 'must be positive';

ALTER TABLE "entry" ADD FOREIGN KEY ("accountId") REFERENCES "account" ("id");

ALTER TABLE "transfer" ADD FOREIGN KEY ("fromAccountId") REFERENCES "account" ("id");

ALTER TABLE "transfer" ADD FOREIGN KEY ("toAccountId") REFERENCES "account" ("id");