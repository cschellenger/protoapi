.PHONY: all clean protobufs

all: protobufs gobuild

protobufs:
	protoc -I. --go_out=. proto/search.proto proto/la_crime.proto
#	protoc -I. -I${GOMODCACHE}/github.com/srikrsna/protoc-gen-gotag\@v1.0.2 --go_out=. proto/search.proto proto/la_crime.proto
#	protoc -I. --gotag_out=. proto/search.proto proto/la_crime.proto
#	protoc -I. -I${GOMODCACHE}/github.com/srikrsna/protoc-gen-gotag\@v1.0.2 --gotag_out=. proto/search.proto proto/la_crime.proto

gobuild: protoapi

protoapi:
	go build .
