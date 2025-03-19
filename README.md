# Protobuf API

Given some simple entity definitions in protobuf format, this will automatically generate a REST API capable of interfacing with a SQL backend.

## Prerequisites

- Install protoc and protoc-gen-go. See [gRPC quickstart](https://grpc.io/docs/languages/go/quickstart/).
- Install [duckdb](https://duckdb.org/#quickinstall).

## Building the DB
- Download [LA crime data CSV](https://data.lacity.org/api/views/2nrs-mtv8/rows.csv?accessType=DOWNLOAD) and replace the headers with the ones from proto/la_crime_headers.csv.
- Import the CSV into duckdb:
```
$ duckdb
D .open open_data.db
D create table la_crime as select * from 'Crime_Data_from_2020_to_Present.csv';
D
```
