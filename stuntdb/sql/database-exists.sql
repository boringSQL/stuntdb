-- name: database-exists
select exists (select 1 from pg_database where datname = $1)
