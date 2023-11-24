CREATE TABLE "user" (
  "username" varchar PRIMARY KEY,
  "name1" varchar NOT NULL,
  "name2" varchar,
  "lastname1" varchar NOT NULL,
  "lastname2" varchar,
  "email" varchar UNIQUE NOT NULL,
  "hashedPassword" varchar NOT NULL,
  "passwordChangedAt" timestamptz NOT NULL DEFAULT (now()),
  "createdAt" timestamptz
);

CREATE TABLE "account" (
  "id" BIGSERIAL PRIMARY KEY,
  "username" varchar NOT NULL,
  "currency" varchar NOT NULL,
  "balance" bigint NOT NULL,
  "createdAt" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "entry" (
  "id" BIGSERIAL PRIMARY KEY,
  "accountId" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "createdAt" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transfer" (
  "id" BIGSERIAL PRIMARY KEY,
  "fromAccountId" bigint NOT NULL,
  "toAccountId" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "createdAt" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "user" ("email");

CREATE INDEX ON "user" ("passwordChangedAt");

CREATE INDEX ON "account" ("username");

CREATE UNIQUE INDEX ON "account" ("username", "currency");

CREATE INDEX ON "entry" ("accountId");

CREATE INDEX ON "entry" ("accountId", "createdAt");

CREATE INDEX ON "transfer" ("fromAccountId");

CREATE INDEX ON "transfer" ("toAccountId");

CREATE INDEX ON "transfer" ("fromAccountId", "toAccountId");

CREATE INDEX ON "transfer" ("toAccountId", "createdAt");

COMMENT ON COLUMN "entry"."amount" IS 'can be positive or negative';

COMMENT ON COLUMN "transfer"."amount" IS 'must be positive';

ALTER TABLE "account" ADD FOREIGN KEY ("username") REFERENCES "user" ("username");

ALTER TABLE "entry" ADD FOREIGN KEY ("accountId") REFERENCES "account" ("id");

ALTER TABLE "transfer" ADD FOREIGN KEY ("fromAccountId") REFERENCES "account" ("id");

ALTER TABLE "transfer" ADD FOREIGN KEY ("toAccountId") REFERENCES "account" ("id");
