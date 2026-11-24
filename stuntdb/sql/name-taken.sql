-- name: name-taken
select exists (select 1 from pg_database where datname = $1)
    or exists (select 1 from pg_roles where rolname = $1)
