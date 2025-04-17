# Protobuf API

Given some simple entity definitions in protobuf format, this will automatically generate a REST API capable of interfacing with a SQL backend.

## Prerequisites

- Install protoc and protoc-gen-go. See [gRPC quickstart](https://grpc.io/docs/languages/go/quickstart/).
- Install [duckdb](https://duckdb.org/#quickinstall).

## Building the DB
- Download [LA crime data CSV](https://data.lacity.org/api/views/2nrs-mtv8/rows.csv?accessType=DOWNLOAD)
- Move CSV to protoapi directory (if not already)
- Import the CSV into duckdb:
```bash
make db
```

## Sample Queries

### Stolen vehicles

The crime code for stolen vehicles is '510'.

#### CEL Query

To represent this as a CEL query, we would write `record.crime_code == 510` and URL encode it. 

```javascript
encodeURIComponent("record.crime_code == 510")
```

#### URL with CEL Query encoded
http://localhost:8080/search?Query=record.crime_code%20%3D%3D%20510
