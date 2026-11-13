#!/bin/bash
# build the lab: a postgres:18 container with pg_regresql installed and
# production-sized data. usage: [LAB_PORT=5432] testdata/lab/up.sh
set -euo pipefail

here=$(cd "$(dirname "$0")" && pwd)
container=${LAB_CONTAINER:-stuntdb-lab}
image=${LAB_IMAGE:-stuntdb-lab}
port=${LAB_PORT:-5432}
pgext=${LAB_PGEXT:-$here/../../../regresql/pg_ext}

if [ ! -d "$pgext" ]; then
	echo "missing $pgext: clone boringsql/regresql beside this repo, or set LAB_PGEXT" >&2
	exit 1
fi

docker build -q -t "$image" -f "$here/Dockerfile" "$pgext/.." >/dev/null
docker rm -f "$container" >/dev/null 2>&1 || true
docker run -d --name "$container" -e POSTGRES_PASSWORD=lab -p "$port:5432" "$image" >/dev/null

ready() { docker exec "$container" psql -U postgres -qc 'select 1' >/dev/null 2>&1; }

# the image restarts its server after initdb
for i in $(seq 1 60); do ready && break; sleep 1; done
sleep 1
for i in $(seq 1 30); do ready && break; sleep 1; done
ready || { docker logs "$container" 2>&1 | tail -5; exit 1; }

docker exec "$container" psql -U postgres -q -c 'create database lab'
docker exec -i "$container" psql -U postgres -q -o /dev/null -v ON_ERROR_STOP=1 -d lab < "$here/schema.sql"
docker exec -i "$container" psql -U postgres -q -o /dev/null -v ON_ERROR_STOP=1 -d lab < "$here/seed.sql"

docker exec "$container" psql -U postgres -At -d lab -c \
	"select format('%s orders, %s items, %s customers', (select count(*) from orders), (select count(*) from order_items), (select count(*) from customers))"
echo "postgres://postgres:lab@localhost:$port/lab"
