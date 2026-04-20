package models

import (
	"database/sql"
	"forum/src/utils"
	"os"
	"path"
	"runtime"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const migrations_dir = "migrations"

var DB *sql.DB

func InitDB(dbPath string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		(&Error{}).Consume(err).LogError()
		return err
	}
	db.Exec("PRAGMA journal_mode=WAL;")
	db.SetMaxOpenConns(1)
	runMigrations(db)
	DB = db
	return nil
}

func runMigrations(db *sql.DB) {
	_, x, _, ok := runtime.Caller(0)
	if !ok {
		(&Error{}).Consume(fmt.Errorf("failed to get caller information")).LogError()
		return
	}
	migrations_dir := path.Join(path.Dir(x), "..", "..", migrations_dir)
	utils.LogDebug(x)
	migrations_found, err := os.ReadDir(migrations_dir)
	if err != nil {
		(&Error{}).Consume(err).LogError()
		return
	}
	for _, file := range migrations_found {
		bytes, err := os.ReadFile(path.Join(migrations_dir, file.Name()))
		if err != nil {
			(&Error{}).Consume(err).LogError()
			return
		}
		query := string(bytes)
		_, err = db.Exec(query)
		if err != nil {
			(&Error{}).Consume(err).LogError()
			return
		}
	}
}
