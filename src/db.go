package forum

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func Init() error {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return err
	}
	defer db.Close()
	var sqliteVersion string
	err = db.QueryRow("select sqlite_version()").Scan(&sqliteVersion)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "%s", sqliteVersion)
	return nil
}
