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
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		panic("DATABASE_URL must be set!")
	}

	sqlFiles := os.Getenv("DB_MIGRATIONS")
	if sqlFiles == "" {
		panic("DB_MIGRATIONS must be set!")
	}

	sqlFileUrl := fmt.Sprintf("file://%s", sqlFiles)
	m, err := migrate.New(sqlFileUrl, databaseUrl)
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}

	v, _, _ := m.Version()
	if v > 0 {
		err = m.Down()
		if err != nil {
			panic(fmt.Sprintf("Failed to reset database: %+v", err))
		}
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		panic(fmt.Sprintf("Failed to complete database migration: %+v", err))
	}

	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		panic(fmt.Sprintf("Unable to open DB connection: %+v", err))
	}

	return &service.Database{DB: db}
}
