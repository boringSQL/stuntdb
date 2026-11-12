ARG PG_MAJOR=18
FROM postgres:${PG_MAJOR}
ARG PG_MAJOR
RUN apt-get update && apt-get install -y --no-install-recommends build-essential postgresql-server-dev-${PG_MAJOR} && rm -rf /var/lib/apt/lists/*
COPY pg_ext /src/pg_ext
RUN cd /src/pg_ext && make clean >/dev/null 2>&1; make && make install
