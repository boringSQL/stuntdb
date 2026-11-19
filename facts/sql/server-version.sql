-- name: server-version
select current_setting('server_version'), current_setting('server_version_num')::int
