.PHONY: build lab test

LAB_PGEXT ?= ../regresql/pg_ext
LAB_IMAGE ?= stuntdb-lab
LAB_CONTAINER ?= stuntdb-lab
LAB_PORT ?= 5432

build:
	go build -o ~/bin/stuntdb ./cmd/stuntdb

test:
	go test ./...

lab:
	@LAB_PGEXT=$(LAB_PGEXT) LAB_IMAGE=$(LAB_IMAGE) LAB_CONTAINER=$(LAB_CONTAINER) LAB_PORT=$(LAB_PORT) testdata/lab/up.sh
