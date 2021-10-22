package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/CircleCI-Public/circleci-demo-docker/service"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	db := SetupDB()
	server := service.NewServer(db)
	http.HandleFunc("/", server.ServeHTTP)
	http.ListenAndServe(":8080", nil)
}

func SetupDB() *service.Database {
	databaseUrl := os.Getenv("CONTACTS_DB_URL")
	if databaseUrl == "" {
		panic("CONTACTS_DB_URL must be set!")
	}

	sqlFiles := os.Getenv("DB_MIGRATIONS")
	if sqlFiles == "" {
		panic("DB_MIGRATIONS must be set!")
	}

	sqlFileUrl := fmt.Sprintf("file://%s", sqlFiles)
	_, err := migrate.New(sqlFileUrl, databaseUrl)
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}

	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		panic(fmt.Sprintf("Unable to open DB connection: %+v", err))
	}

	return &service.Database{DB: db}
}
