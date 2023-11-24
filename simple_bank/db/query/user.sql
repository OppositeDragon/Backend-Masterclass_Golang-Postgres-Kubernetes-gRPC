-- name: CreateUser :one 
INSERT INTO "user" (
    username,
    name1,
    name2,
    lastname1,
    lastname2,
    email,
    "hashedPassword"
  )
VALUES($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetUser :one
SELECT username,
  name1,
  name2,
  lastname1,
  lastname2,
  email,
  "hashedPassword",
  "passwordChangedAt",
  "createdAt"
FROM "user"
WHERE username = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT username,
  name1,
  name2,
  lastname1,
  lastname2,
  email,
  "hashedPassword",
  "passwordChangedAt",
  "createdAt"
FROM "user"
WHERE email = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE "user"
SET
  name1 = $2,
  name2 = $3,
  lastname1 = $4,
  lastname2 = $5,
  email = $6,
  "hashedPassword" = $7,
  "passwordChangedAt" = $8
WHERE username = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM "user"
WHERE username = $1;

-- name: GetUsers :many
SELECT username,
  name1,
  name2,
  lastname1,
  lastname2,
  email,
  "hashedPassword",
  "passwordChangedAt",
  "createdAt"
FROM "user"
ORDER BY username
LIMIT $1 OFFSET $2;