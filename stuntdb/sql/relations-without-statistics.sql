-- name: relations-without-statistics
select n.nspname, c.relname
from pg_class c
join pg_namespace n on n.oid = c.relnamespace
where c.relkind in ('r', 'i', 'm')
  and n.nspname not in ('pg_catalog', 'information_schema')
  and n.nspname !~ '^pg_'
  and c.reltuples < 0
order by n.nspname, c.relname
