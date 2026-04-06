package postgres

import (
	"database/sql"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var database *sql.DB

func GetDb() *sql.DB {
	if database != nil {
		return database
	}
	_ = openConnection()
	return database
}

func Close() error {
	return database.Close()
}

func openConnection() error {
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	database = db
	return nil
}