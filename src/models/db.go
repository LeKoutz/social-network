package models

import (
	"database/sql"
	"errors"
	utils "forum/src/utils"

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

func GetAllUsernames() ([]string, error) {
	rows, err := DB.Query(`SELECT username FROM users`)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return []string{}, err
	}
	defer rows.Close()
	var usernames []string
	for rows.Next() {
		var email string
		err = rows.Scan(&email)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return []string{}, err
		}
		usernames = append(usernames, email)
	}
	return usernames, nil
}

func GetAllUserEmails() ([]string, error) {
	rows, err := DB.Query(`SELECT email FROM users`)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return []string{}, err
	}
	defer rows.Close()
	var emails []string
	for rows.Next() {
		var email string
		err = rows.Scan(&email)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return []string{}, err
		}
		emails = append(emails, email)
	}
	return emails, nil
}

func SetUserSession(id int64, session_key string) error {
	stmt, err := DB.Prepare("UPDATE users SET session_key = ? WHERE id = ?")
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	_, err = stmt.Exec(session_key, id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	return nil
}

func GetLikedPostsByUserId(user_id int64) (Posts, error) {
	var posts Posts
	rows, err := DB.Query(`
	SELECT posts.id, posts.title, posts.body, posts.timestamp
	FROM posts
	JOIN reactions r ON posts.id = r.post_id
	WHERE r.user_id = ? AND r.value = 1
	`, user_id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return Posts{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var post Post
		var ts string
		err = rows.Scan(&post.Id, &post.Title, &post.Body, &ts)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return Posts{}, err
		}
		t, err := utils.ConvertStringToTime(ts)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return Posts{}, err
		}
		post.TimestampString = utils.ConvertTimeToString(t)
		posts = append(posts, post)
	}
	return posts, nil
}

func GetPostsByUserId(id int64) (Posts, error) {
	var posts Posts
	rows, err := DB.Query(`
	SELECT posts.id, posts.title, posts.body, posts.timestamp
	FROM posts
	WHERE user_id = ?`, id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return Posts{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var post Post
		var ts string
		err = rows.Scan(&post.Id, &post.Title, &post.Body, &ts)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return Posts{}, err
		}
		t, err := utils.ConvertStringToTime(ts)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return Posts{}, err
		}
		post.TimestampString = utils.ConvertTimeToString(t)
		posts = append(posts, post)
	}
	return posts, nil
}

func AddCategoryToPost(post_id, category_id int64) error {
	stmt, err := DB.Prepare("INSERT INTO posts_categories (post_id, category_id) VALUES (?, ?)")
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	_, err = stmt.Exec(post_id, category_id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	return nil
}

func GetCommentsByPostId(postId int64) (Comments, error) {
	rows, err := DB.Query(`
	SELECT
	c.id,
	c.post_id,
	c.user_id,
	c.body,
	c.timestamp,
	u.username
	FROM comments c
	JOIN users u ON c.user_id = u.id
	WHERE c.post_id = ?
	ORDER BY c.timestamp ASC`, postId)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return Comments{}, err
	}
	defer rows.Close()
	var comments Comments
	for rows.Next() {
		var comment Comment
		var ts string
		err = rows.Scan(
			&comment.Id,
			&comment.PostId,
			&comment.UserId,
			&comment.Body,
			&ts,
			&comment.Username,
		)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return Comments{}, err
		}
		t, err := utils.ConvertStringToTime(ts)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return Comments{}, err
		}
		comment.TimestampString = utils.ConvertTimeToString(t)
		comments = append(comments, comment)
	}
	return comments, nil
}
