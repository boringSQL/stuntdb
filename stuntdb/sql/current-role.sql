-- name: current-role
select current_user, rolsuper
from pg_roles
where rolname = current_user
