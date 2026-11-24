-- name: set-autovacuum
alter table %s set (autovacuum_enabled = false, toast.autovacuum_enabled = false)
