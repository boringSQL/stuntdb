-- name: catalog-counts
select
  (select count(*)
     from pg_class c join pg_namespace n on n.oid = c.relnamespace
    where c.relkind in ('r', 'p', 'm') and n.nspname not in ('pg_catalog', 'information_schema')
      and n.nspname !~ '^pg_'),
  (select count(*)
     from pg_class c join pg_namespace n on n.oid = c.relnamespace
    where c.relkind = 'i' and n.nspname not in ('pg_catalog', 'information_schema')
      and n.nspname !~ '^pg_'),
  (select count(*)
     from pg_attribute a
     join pg_class c on c.oid = a.attrelid
     join pg_namespace n on n.oid = c.relnamespace
    where c.relkind in ('r', 'p', 'm') and a.attnum > 0 and not a.attisdropped
      and n.nspname not in ('pg_catalog', 'information_schema') and n.nspname !~ '^pg_')
