-- name: CreateEntry :one
insert into entry("accountId", amount)
values($1, $2)
RETURNING *;

-- name: GetEntry :one
select *
from entry
where id = $1
limit 1;

-- name: GetEntries :many
select *
from entry
order by id
limit $1 offset $2;

-- name: UpdateEntry :one
update entry
set amount = $2
where id = $1
RETURNING *;

-- name: DeleteEntry :exec
delete from entry
where id = $1;