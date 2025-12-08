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
	stmt, err := db.Prepare("INSERT INTO users (username, email, hash) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(user.Username, user.Email, user.Hash)
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
	err = db.QueryRow(`SELECT id, email, username, hash FROM users WHERE email = ?`, email).Scan(&user.Id, &user.Email, &user.Username, &user.Hash)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func getUserBySession(sessionValue string) (User, error) {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return User{}, err
	}
	defer db.Close()
	var user User
	err = db.QueryRow(`SELECT id, email, username, hash FROM users WHERE session_key = ?`, sessionValue).Scan(&user.Id, &user.Email, &user.Username, &user.Hash)
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

func setUserSession(id int, session_key string) error {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare("UPDATE users SET session_key = ? WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(session_key, id)
	if err != nil {
		return err
	}
	return nil
}

func getAllCategories() (Categories, error) {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return []Category{}, err
	}
	defer db.Close()
	rows, err := db.Query(`SELECT id, name FROM categories`)
	if err != nil {
		return []Category{}, err
	}
	defer rows.Close()
	var categories Categories
	for rows.Next() {
		var category Category
		err = rows.Scan(&category.Id, &category.Name)
		if err != nil {
			return []Category{}, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func getCategoryById(id int) (Category, error) {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return Category{}, err
	}
	defer db.Close()
	var category Category
	err = db.QueryRow(`SELECT name FROM categories WHERE id = ?`, id).Scan(&category.Name)
	if err != nil {
		return Category{}, err
	}
	category.Id = id
	return category, nil
}

func addCategory(category Category) error {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return err
	}
	defer db.Close()
	err = category.validateCategory()
	if err != nil {
		return err
	}
	if (&category).DoesCategoryExist() {
		return ErrorCategoryAlreadyExists
	}
	stmt, err := db.Prepare("INSERT INTO categories (name) VALUES (?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(category.Name)
	if err != nil {
		return err
	}
	return nil
}

func getPostsByCategoryId(id int) (Posts, error) {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return Posts{}, err
	}
	defer db.Close()
	var posts Posts
	rows, err := db.Query(`
	SELECT posts.id, posts.title, posts.body
	FROM posts
	JOIN posts_categories pc ON posts.id = pc.post_id
	JOIN categories ON pc.category_id = categories.id
	WHERE pc.category_id = ?`, id)
	if err != nil {
		return Posts{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var post Post
		rows.Scan(&post.Id, &post.Title, &post.Body)
		posts = append(posts, post)
	}
	return posts, nil
}

func getAllPosts() (Posts, error) {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return Posts{}, err
	}
	defer db.Close()
	rows, err := db.Query(`SELECT id, title, body FROM posts`)
	if err != nil {
		return Posts{}, err
	}
	defer rows.Close()
	var posts Posts
	for rows.Next() {
		var post Post
		err = rows.Scan(&post.Id, &post.Title, &post.Body)
		if err != nil {
			return Posts{}, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func getPostById(id int) (Post, error) {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return Post{}, err
	}
	defer db.Close()
	var post Post
	err = db.QueryRow(`SELECT title, body FROM posts WHERE id = ?`, id).Scan(&post.Title, &post.Body)
	if err != nil {
		return Post{}, err
	}
	post.Id = id
	return post, nil
}

func addPost(post Post) (int, error) {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return 0, err
	}
	defer db.Close()
	err = post.validatePost()
	if err != nil {
		return 0, err
	}
	stmt, err := db.Prepare("INSERT INTO posts (title, body, user_id) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(post.Title, post.Body, post.UserId)
	if err != nil {
		return 0, err
	}

	// Get the last inserted post ID
	postId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Insert into posts_categories table
	stmt, err = db.Prepare("INSERT INTO posts_categories (post_id, category_id) VALUES (?, ?)")
	if err != nil {
		return 0, err
	}
	_, err = stmt.Exec(postId, post.Category.Id)
	if err != nil {
		return 0, err
	}

	return int(postId), nil
}

func addComment(comment Comment) (int, error) {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return 0, err
	}
	defer db.Close()

	if err := comment.validateComment(); err != nil {
		return 0, err
	}

	res, err := db.Exec(
		"INSERT INTO comments (post_id, user_id, body) VALUES (?, ?, ?)",
		comment.PostId,
		comment.UserId,
		comment.Body,
	)
	commentId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(commentId), nil
}

func getCommentsByPostId(postId int) (Comments, error) {
	db, err := sql.Open("sqlite3", "./db.db:")
	if err != nil {
		return Comments{}, err
	}
	defer db.Close()

	rows, err := db.Query(`
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
		return Comments{}, err
	}
	defer rows.Close()

	var comments Comments
	for rows.Next() {
		var comment Comment

		err = rows.Scan(
			&comment.Id,
			&comment.PostId,
			&comment.UserId,
			&comment.Body,
			&comment.Timestamp,
			&comment.Username,
		)
		if err != nil {
			return Comments{}, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func getReactionsForPost(postId int) (int, int, error) {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return 0, 0, err
	}
	defer db.Close()

	var likes int
	err = db.QueryRow(`
        SELECT COUNT(*)
        FROM reactions
        WHERE post_id = ? AND value = 1
    `, postId).Scan(&likes)
	if err != nil {
		return 0, 0, err
	}

	var dislikes int
	err = db.QueryRow(`
        SELECT COUNT(*)
        FROM reactions
        WHERE post_id = ? AND value = 2
    `, postId).Scan(&dislikes)
	if err != nil {
		return 0, 0, err
	}

	return likes, dislikes, nil
}
