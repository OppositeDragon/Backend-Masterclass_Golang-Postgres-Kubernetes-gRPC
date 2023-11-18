-- name: CreateTransfer :one
insert into transfer("fromAccountId", "toAccountId", "amount")
values($1, $2, $3)
RETURNING *;

-- name: GetTransfer :one
select *
from transfer
where id = $1
limit 1;

-- name: GetTransfers :many
select *
from transfer
order by id
limit $1 offset $2;

-- name: UpdateTransfer :one
update transfer
set amount = $2
where id = $1
RETURNING *;

-- name: DeleteTransfer :exec
delete from transfer
where id = $1;