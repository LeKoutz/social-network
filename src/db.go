package forum

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func Init() error {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec(createUsersTable())
	if err != nil {
		return err
	}
	_, err = db.Exec(createPostsTable())
	if err != nil {
		return err
	}
	_, err = db.Exec(createCategoriesTable())
	if err != nil {
		return err
	}
	_, err = db.Exec(createCommentsTable())
	if err != nil {
		return err
	}
	_, err = db.Exec(createReactionsTable())
	if err != nil {
		return err
	}
	_, err = db.Exec(createPostsCategoriesTable())
	if err != nil {
		return err
	}
	return nil
}

func createUsersTable() string {
	return `CREATE TABLE IF NOT EXISTS "users" (
		"id"	INTEGER NOT NULL UNIQUE,
		"email"	TEXT NOT NULL UNIQUE,
		"username"	TEXT NOT NULL,
		"salt"	TEXT UNIQUE,
		"hash"	TEXT,
		"session_key"	TEXT,
		PRIMARY KEY("id" AUTOINCREMENT)
	)`
}

func createPostsTable() string {
	return `CREATE TABLE IF NOT EXISTS "posts" (
		"id"	INTEGER NOT NULL UNIQUE,
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
		"name"	TEXT,
		PRIMARY KEY("id" AUTOINCREMENT)
	)`
}

func createCommentsTable() string {
	return `CREATE TABLE IF NOT EXISTS "comments" (
		"id"	INTEGER NOT NULL UNIQUE,
		"post_id"	INTEGER NOT NULL,
		"user_id"	INTEGER NOT NULL,
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

func registerUserOnDB(user User) error {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return err
	}
	defer db.Close()
	err = user.validateUser()
	if err != nil {
		return err
	}
	if IsEmailRegistered(user.Email) {
		return ErrorEmailIsRegistered
	}
	stmt, err := db.Prepare("INSERT INTO users (username, email, salt, hash) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(user.Username, user.Email, user.Salt, user.Hash)
	if err != nil {
		return err
	}
	return nil
}

func getAllUsernames() ([]string, error) {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return []string{}, err
	}
	defer db.Close()
	rows, err := db.Query(`SELECT username FROM users`)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()
	var usernames []string
	for rows.Next() {
		var email string
		err = rows.Scan(&email)
		if err != nil {
			return []string{}, err
		}
		usernames = append(usernames, email)
	}
	return usernames, nil
}

func getAllUserEmails() ([]string, error) {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return []string{}, err
	}
	defer db.Close()
	rows, err := db.Query(`SELECT email FROM users`)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()
	var emails []string
	for rows.Next() {
		var email string
		err = rows.Scan(&email)
		if err != nil {
			return []string{}, err
		}
		emails = append(emails, email)
	}
	return emails, nil
}

func getUserByEmail(email string) (User, error) {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return User{}, err
	}
	defer db.Close()
	var user User
	err = db.QueryRow(`SELECT email FROM users`).Scan(&user.Email)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func checkIfUserAlreadyLikedPost(userId, postId int) (int, error) {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var existingReactionId int
	err = db.QueryRow(`
		SELECT id FROM reactions
		WHERE user_id = ? AND post_id = ? AND value = 1
		`, userId, postId).Scan(&existingReactionId)

	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	return existingReactionId, nil
}

func checkIfUserAlreadyDislikedPost(userId, postId int) (int, error) {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var existingDislikeId int
	err = db.QueryRow(`
		SELECT id FROM reactions
		WHERE user_id = ? AND post_id = ? AND value = 2
		`, userId, postId).Scan(&existingDislikeId)

	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	return existingDislikeId, nil
}

func addLikeToPost(userId, postId int) error {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
		INSERT INTO reactions (user_id, post_id, value)
		VALUES (?, ?, 1)
		`, userId, postId)
	if err != nil {
		return err
	}

	return nil
}

func removeDislikeFromPost(dislikeId int) error {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
		DELETE FROM reactions
		WHERE id = ?
		`, dislikeId)
	if err != nil {
		return err
	}

	return nil
}

func addDislikeToPost(userId, postId int) error {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
		INSERT INTO reactions (user_id, post_id, value)
		VALUES (?, ?, 2)
		`, userId, postId)
	if err != nil {
		return err
	}

	return nil
}

func removeLikeFromPost(userId, postId int) error {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
		DELETE FROM reactions
		WHERE user_id = ? AND post_id = ? AND value = 1
		`, userId, postId)
	if err != nil {
		return err
	}

	return nil
}

func checkIfUserAlreadyLikedComment(userId, commentId int) (int, error) {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var existingReactionId int
	err = db.QueryRow(`
		SELECT id FROM reactions
		WHERE user_id = ? AND comment_id = ? AND value = 1
		`, userId, commentId).Scan(&existingReactionId)

	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	return existingReactionId, nil
}

func checkIfUserAlreadyDislikedComment(userId, commentId int) (int, error) {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var existingDislikeId int
	err = db.QueryRow(`
		SELECT id FROM reactions
		WHERE user_id = ? AND comment_id = ? AND value = 2
		`, userId, commentId).Scan(&existingDislikeId)

	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	return existingDislikeId, nil
}

func addLikeToComment(userId, commentId int) error {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
		INSERT INTO reactions (user_id, comment_id, value)
		VALUES (?, ?, 1)
		`, userId, commentId)
	if err != nil {
		return err
	}

	return nil
}

func addDislikeToComment(userId, commentId int) error {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
		INSERT INTO reactions (user_id, comment_id, value)
		VALUES (?, ?, 2)
		`, userId, commentId)
	if err != nil {
		return err
	}

	return nil
}
func removeDislikeFromComment(dislikeId int) error {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
		DELETE FROM reactions
		WHERE id = ?
		`, dislikeId)
	if err != nil {
		return err
	}

	return nil
}

func removeLikeFromComment(userId, commentId int) error {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
		DELETE FROM reactions
		WHERE user_id = ? AND comment_id = ? AND value = 1
		`, userId, commentId)
	if err != nil {
		return err
	}

	return nil
}
