package main

import (
	"darkport.net/protoapi/api"
	"darkport.net/protoapi/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/marcboeker/go-duckdb"
)

func main() {

	db, err := sqlx.Open("duckdb", "open_data.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	server := api.NewSqlServer(db, (&model.CrimeData{}).ProtoReflect(), "la_crime")
	defer server.Close()

	server.Serve()
}
