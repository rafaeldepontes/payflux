package postgres

import (
	"database/sql"
	"os"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"
)

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
		db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
		if err != nil {
			return
		}
		database = db
	})
	return nil
}
