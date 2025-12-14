package main

import (
	forum "forum/src"
	"os"
)

func main() {
	db_path := "./db.db"
	if len(os.Args) != 2 {
		db_path = os.Args[1]
	}
	forum.MockGen(db_path)
}
