CREATE TABLE "session" (
  "id" uuid PRIMARY KEY,
  "username" varchar NOT NULL,
  "access_token" varchar NOT NULL,
  "access_expires_at" timestamptz NOT NULL,
  "refresh_token" varchar NOT NULL,
  "refresh_expires_at" timestamptz NOT NULL,
  "user_agent" varchar,
  "client_ip" varchar,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "createdAt" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "session" ADD FOREIGN KEY ("username") REFERENCES "user"("username");