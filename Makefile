.PHONY: all clean protobufs db

all: protobufs gobuild

db:
	scripts/fix_csv_headers.sh
	duckdb <scripts/create_db.sql

protobufs:
	protoc -I. --go_out=. proto/search.proto proto/la_crime.proto

gobuild: protoapi

protoapi:
	go build .

test:
	go test -v ./...
