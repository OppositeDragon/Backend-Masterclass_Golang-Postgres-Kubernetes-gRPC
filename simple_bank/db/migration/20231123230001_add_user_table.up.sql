ALTER TABLE "account" RENAME COLUMN "owner" TO "username";

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


CREATE INDEX ON "user" ("email");
CREATE INDEX ON "user" ("passwordChangedAt");

ALTER TABLE "account" ADD FOREIGN KEY ("username") REFERENCES "user" ("username");
CREATE UNIQUE INDEX ON "account" ("username", "currency");
