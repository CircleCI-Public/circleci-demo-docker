package test

import (
	"database/sql"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/CircleCI-Public/circleci-demo-docker/service"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
)

// Env provides access to all services used in tests, like the database, our server, and an HTTP client for performing
// HTTP requests against the test server.
type Env struct {
	T          *testing.T
	DB         *service.Database
	Server     *service.Server
	HttpServer *httptest.Server
	Client     service.Client
}

// Close must be called after each test to ensure the Env is properly destroyed.
func (env *Env) Close() {
	env.HttpServer.Close()
	env.DB.Close()
}

// SetupEnv creates a new test environment, including a clean database and an instance of our HTTP service.
func SetupEnv(t *testing.T) *Env {
	db := SetupDB(t)
	server := service.NewServer(db)
	httpServer := httptest.NewServer(server)
	return &Env{
		T:          t,
		DB:         db,
		Server:     server,
		HttpServer: httpServer,
		Client:     service.NewClient(httpServer.URL),
	}
}

// SetupDB initializes a test database, performing all migrations.
func SetupDB(t *testing.T) *service.Database {
	databaseUrl := os.Getenv("DATABASE_URL")
	require.NotEmpty(t, databaseUrl, "DATABASE_URL must be set!")

	sqlFiles := os.Getenv("DB_MIGRATIONS")
	require.NotEmpty(t, sqlFiles, "DB_MIGRATIONS must be set!")

	sqlFileUrl := fmt.Sprintf("file://%s", sqlFiles)
	m, err := migrate.New(sqlFileUrl, databaseUrl)
	require.NoError(t, err, "Failed to start database migration")
	v, _, _ := m.Version()
	if v > 0 {
		err = m.Down()
		require.NoError(t, err, "Failed to reset database")
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		require.Fail(t, "Failed to complete database migration")
	}

	db, err := sql.Open("postgres", databaseUrl)
	require.NoError(t, err, "Error opening database")

	return &service.Database{db}
}
