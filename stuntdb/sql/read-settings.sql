-- name: read-settings
select name, setting
from pg_settings
where name = any($1)
