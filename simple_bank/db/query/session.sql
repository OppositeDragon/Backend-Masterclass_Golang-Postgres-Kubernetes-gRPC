-- name: CreateSession :one 
INSERT INTO "session" (
    id,
    username,
    access_token,
    access_expires_at,
    refresh_token,
    refresh_expires_at,
    user_agent,
    client_ip
  )
VALUES(
  $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: UpdateSessionAccess :one
UPDATE "session"
SET
  access_token = $2,
  access_expires_at = $3
WHERE "id" = $1;

-- name: UpdateSessionRefresh :one
UPDATE "session"
SET
  access_token = $2,
  access_expires_at = $3
  refresh_token = $4,
  refresh_expires_at = $5
WHERE "id" = $1;

-- name: GetSession :one
SELECT * FROM "session"
WHERE "id" = $1
LIMIT 1;

-- name: DeleteSession :exec
DELETE FROM "session"
WHERE "id" = $1;