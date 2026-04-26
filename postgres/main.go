package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"time"

	_ "github.com/lib/pq" // postgres driver

	"github.com/swarm-deploy/init-jobs/postgres/internal"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)

	exit := func(code int) {
		cancel()
		os.Exit(code)
	}

	dsnString, err := readDSNSecret()
	if err != nil {
		slog.Error("You must provide secret in file /run/secrets/dsn")
		exit(1)
	}

	dsn, err := internal.ParseDSN(dsnString)
	if err != nil {
		slog.Error("failed to parse DSN", slog.Any("err", err))
		exit(1)
	}

	slog.Info("dsn found", slog.String("dbname", dsn.DatabaseName))

	if valid := validateDatabaseName(dsn.DatabaseName); !valid {
		slog.Error("database name invalid", slog.String("database_name", dsn.DatabaseName))
		exit(1)

		return
	}

	db, err := sql.Open("postgres", dsn.PostgresConnectionString)
	if err != nil {
		slog.Error("failed to open database connection", slog.Any("err", err))
		exit(1)
	}

	alreadyExists, err := checkDatabaseExists(ctx, db, dsn.DatabaseName)
	if err != nil {
		slog.Error("failed to check database exists", slog.Any("err", err))
		exit(1)
	}

	if alreadyExists {
		slog.Info("database already exists")
		exit(0)
	}

	query := fmt.Sprintf("CREATE DATABASE %s", dsn.DatabaseName)

	if _, err = db.ExecContext(ctx, query); err != nil {
		slog.Error("failed to execute query", slog.Any("err", err))
		exit(1)
	}

	slog.Info("database created")
}

func readDSNSecret() (string, error) {
	data, err := os.ReadFile("/run/secrets/dsn")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func checkDatabaseExists(ctx context.Context, db *sql.DB, databaseName string) (bool, error) {
	var exists bool

	query := "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)"
	err := db.QueryRowContext(ctx, query, databaseName).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func validateDatabaseName(databaseName string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !re.MatchString(databaseName) {
		return false
	}

	return true
}
