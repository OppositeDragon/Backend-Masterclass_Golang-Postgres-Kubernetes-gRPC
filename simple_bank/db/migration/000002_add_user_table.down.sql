alter table if exists "account" drop constraint if exists "account_username_fkey";
alter table if exists "account" drop constraint if exists "username_currency_key";

drop table if exists "user";
