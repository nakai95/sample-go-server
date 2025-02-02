package test

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go"
	postgresContainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var db *sql.DB

func SetupPostgresContainer(t *testing.T) (*sql.DB, error) {
	t.Helper()
	ctx := context.Background()

	dbName := "sample"
	dbUser := "user"
	dbPassword := "password"

	pgContainer, err := postgresContainer.Run(ctx,
		"postgres:latest",
		postgresContainer.WithDatabase(dbName),
		postgresContainer.WithUsername(dbUser),
		postgresContainer.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)

	if err != nil {
		t.Fatalf("failed to start container: %s", err)
	}

	t.Cleanup(func() {
		t.Logf("terminating container")
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	// Get port
	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("failed to get mapped port: %s", err)
	}

	// Run migrations
	if err := runMigrations(dbName, dbUser, dbPassword, port.Int()); err != nil {
		t.Fatalf("failed to run migrations: %s", err)
	}

	// Get connection to the database
	databaseURL := fmt.Sprintf("postgres://%s:%s@localhost:%d/%s?sslmode=disable", dbUser, dbPassword, port.Int(), dbName)
	db, err = sql.Open("postgres", databaseURL)
	if err != nil {
		t.Fatalf("failed to open database: %s", err)
	}
	return db, nil
}

func runMigrations(dbName string, dbUser string, dbPassword string, port int) error {
	databaseURL := fmt.Sprintf("postgres://%s:%s@localhost:%d/%s?sslmode=disable", dbUser, dbPassword, port, dbName)

	projectRoot, err := getProjectRoot()
	if err != nil {
		return fmt.Errorf("failed to get project root: %w", err)
	}

	migrationsFilePath := filepath.Join(projectRoot, "migrate/postgres/migrations")

	m, err := migrate.New(
		"file://"+migrationsFilePath,
		databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}
	if err := m.Up(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func getProjectRoot() (string, error) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	return filepath.Abs(filepath.Join(basepath, ".."))
}
