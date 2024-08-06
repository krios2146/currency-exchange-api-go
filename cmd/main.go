package main

import (
	"github.com/krios2146/currency-exchange-api-go/internal/api"
	"github.com/krios2146/currency-exchange-api-go/internal/db"
)

func main() {
	server := api.NewServer(db.NewSqliteDBConnection())
	server.Run()
}
