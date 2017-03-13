package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/circleci/cci-demo-docker/service"
	_ "github.com/mattes/migrate/driver/postgres"
	"github.com/mattes/migrate/migrate"
)

func main() {
	db := SetupDB()
	server := service.NewServer(db)
	http.HandleFunc("/", server.ServeHTTP)
	http.ListenAndServe(":8080", nil)
}

func SetupDB() *service.Database {
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		panic("DATABASE_URL must be set!")
	}

	allErrors, ok := migrate.ResetSync(databaseUrl, "./db/migrations")
	if !ok {
		panic(fmt.Sprintf("%+v", allErrors))
	}

	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		panic(fmt.Sprintf("Unable to open DB connection: %+v", err))
	}

	return &service.Database{db}
}
