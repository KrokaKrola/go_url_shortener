package store

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

const migrationsPath = "internal/store/migrations"

type Database struct {
	Connection      *pgx.Conn
	migrationsTable string
}

func NewDatabase() (*Database, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_URL"))

	if err != nil {
		return nil, err
	}

	return &Database{
		Connection:      conn,
		migrationsTable: "_migrations",
	}, nil
}

func (d *Database) Migrate() error {
	config, err := pgx.ParseConfig(os.Getenv("DB_URL"))
	if err != nil {
		return err
	}

	db := stdlib.OpenDB(*config)

	driver, err := postgres.WithInstance(db, &postgres.Config{
		MigrationsTable: "_migrations",
	})

	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres", driver)

	if err != nil {
		return err
	}

	version, _, error := m.Version()

	if error != nil {
		return error
	}

	log.Info("Current database version: ", version)

	if err := m.Up(); err != nil {
		if err.Error() == "no change" {
			log.Info("Database is up to date")
			return nil
		}

		return err
	}

	return nil
}
