-- name: CreateAccount :one
insert into account(username, balance, currency)
values($1, $2, $3)
RETURNING *;

-- name: GetAccount :one
select *
from account
where id = $1
limit 1;

-- name: GetAccountForUpdate :one
select *
from account
where id = $1
limit 1
for no key update;

-- name: GetAccounts :many
select *
from account
order by id
limit $1 offset $2;

-- name: UpdateAccount :one
update account
set balance = $2
where id = $1
RETURNING *;

-- name: AddAmountAccount :one
update account
set balance = balance + sqlc.arg(amount)
where id = $1
RETURNING *;

-- name: DeleteAccount :exec
delete from account
where id = $1;