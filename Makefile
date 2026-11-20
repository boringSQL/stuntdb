.PHONY: build facts lab test

LAB_PGEXT ?= ../regresql/pg_ext
LAB_IMAGE ?= stuntdb-lab
LAB_CONTAINER ?= stuntdb-lab
LAB_PORT ?= 5432
LAB_URI ?= postgres://postgres:lab@localhost:$(LAB_PORT)/lab

build:
	go build -o ~/bin/stuntdb ./cmd/stuntdb

facts: build
	rm -rf testdata/lab/facts
	~/bin/stuntdb facts $(LAB_URI) -o testdata/lab/facts

test:
	go test ./...

lab:
	@LAB_PGEXT=$(LAB_PGEXT) LAB_IMAGE=$(LAB_IMAGE) LAB_CONTAINER=$(LAB_CONTAINER) LAB_PORT=$(LAB_PORT) testdata/lab/up.sh
