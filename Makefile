.PHONY: all clean protobufs

all: protobufs gobuild

protobufs:
	protoc -I. --go_out=. proto/search.proto proto/la_crime.proto

gobuild: protoapi

protoapi:
	go build .
