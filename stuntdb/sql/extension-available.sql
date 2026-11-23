-- name: extension-available
select exists (select 1 from pg_available_extensions where name = 'pg_regresql')
