package postgres

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

var (
	database *sql.DB
	once     sync.Once
)

func GetDb() *sql.DB {
	if database != nil {
		return database
	}
	_ = openConnection()
	return database
}

func Close() error {
	if database != nil {
		return database.Close()
	}
	return nil
}

func openConnection() error {
	once.Do(func() {
		dbURL := os.Getenv("DATABASE_URL")
		var db *sql.DB
		var err error

		for i := 0; i < 5; i++ {
			db, err = sql.Open("pgx", dbURL)
			if err == nil {
				err = db.Ping()
			}

			if err == nil {
				database = db
				return
			}

			log.Printf("Waiting for database... (attempt %d/5): %v", i+1, err)
			time.Sleep(2 * time.Second)
		}

		if err != nil {
			log.Fatalf("Could not connect to database after 5 attempts: %v", err)
		}
	})
	return nil
}

func RunMigrations() error {
	db := GetDb()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create migrate driver: %v", err)
	}

	source, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("could not create source driver: %v", err)
	}

	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run up migrations: %v", err)
	}

	log.Println("Migrations applied successfully!")
	return nil
}
