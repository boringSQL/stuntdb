-- name: user-schemas
select nspname
from pg_namespace
where nspname not in ('pg_catalog', 'information_schema')
  and nspname !~ '^pg_'
order by nspname
