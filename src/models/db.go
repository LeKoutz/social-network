package models

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(dbPath string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		(&Error{}).Consume(err).LogError()
		return err
	}
	db.Exec("PRAGMA journal_mode=WAL;")
	db.SetMaxOpenConns(1)
	_, err = db.Exec(createUsersTable())
	if err != nil {
		(&Error{}).Consume(err).LogError()
	}
	_, err = db.Exec(createPostsTable())
	if err != nil {
		(&Error{}).Consume(err).LogError()
	}
	_, err = db.Exec(createCategoriesTable())
	if err != nil {
		(&Error{}).Consume(err).LogError()
	}
	_, err = db.Exec(createCommentsTable())
	if err != nil {
		(&Error{}).Consume(err).LogError()
	}
	_, err = db.Exec(createReactionsTable())
	if err != nil {
		(&Error{}).Consume(err).LogError()
	}
	_, err = db.Exec(createPostsCategoriesTable())
	if err != nil {
		(&Error{}).Consume(err).LogError()
	}
	DB = db
	return nil
}

func createUsersTable() string {
	return `CREATE TABLE IF NOT EXISTS "users" (
		"id"	INTEGER NOT NULL UNIQUE,
		"email"	TEXT NOT NULL UNIQUE,
		"username"	TEXT NOT NULL,
		"hash"	TEXT,
		"session_key"	TEXT,
		PRIMARY KEY("id" AUTOINCREMENT)
	)`
}

func createPostsTable() string {
	return `CREATE TABLE IF NOT EXISTS "posts" (
		"id"	INTEGER NOT NULL UNIQUE,
		"timestamp"	TEXT,
		"title"	TEXT,
		"body"	TEXT,
		"user_id"	INTEGER NOT NULL,
		PRIMARY KEY("id" AUTOINCREMENT),
		FOREIGN KEY("user_id") REFERENCES "users"("id")
	)`
}

func createCategoriesTable() string {
	return `CREATE TABLE IF NOT EXISTS "categories" (
		"id"	INTEGER NOT NULL UNIQUE,
		"name"	TEXT NOT NULL UNIQUE,
		"description"	TEXT NOT NULL UNIQUE,
		PRIMARY KEY("id" AUTOINCREMENT)
	)`
}

func createCommentsTable() string {
	return `CREATE TABLE IF NOT EXISTS "comments" (
		"id"	INTEGER NOT NULL UNIQUE,
		"post_id"	INTEGER NOT NULL,
		"user_id"	INTEGER NOT NULL,
		"timestamp"	TEXT,
		"body"	TEXT NOT NULL,
		PRIMARY KEY("id" AUTOINCREMENT),
		FOREIGN KEY("post_id") REFERENCES "posts"("id"),
		FOREIGN KEY("user_id") REFERENCES "users"("id")
	)`
}

func createPostsCategoriesTable() string {
	return `CREATE TABLE IF NOT EXISTS "posts_categories" (
		"id"	INTEGER NOT NULL UNIQUE,
		"post_id"	INTEGER NOT NULL,
		"category_id"	INTEGER NOT NULL,
		PRIMARY KEY("id" AUTOINCREMENT),
		FOREIGN KEY("category_id") REFERENCES "categories"("id"),
		FOREIGN KEY("post_id") REFERENCES "posts"("id")
	)`
}

func createReactionsTable() string {
	return `CREATE TABLE IF NOT EXISTS "reactions" (
		"id"	INTEGER NOT NULL UNIQUE,
		"post_id"	INTEGER,
		"user_id"	INTEGER NOT NULL,
		"comment_id"	INTEGER,
		"value"	INTEGER NOT NULL,
		PRIMARY KEY("id" AUTOINCREMENT),
		FOREIGN KEY("comment_id") REFERENCES "comments"("id"),
		FOREIGN KEY("post_id") REFERENCES "posts"("id"),
		FOREIGN KEY("user_id") REFERENCES "users"("id")
	)`
}
