-- name: planner-settings
select name, setting
from pg_settings
where context = 'user'
  and (category like 'Query Tuning%' or name = any($1))
order by name
