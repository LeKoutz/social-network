package forum

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init() error {
	db, err := sql.Open("sqlite3", "./db.db")
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
		"name"	TEXT,
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

func registerUserOnDB(user User) error {
	err := user.validateUser()
	if err != nil {
		return err
	}
	if IsEmailRegistered(user.Email) {
		return ErrorEmailIsRegistered
	}
	stmt, err := DB.Prepare("INSERT INTO users (username, email, hash) VALUES (?, ?, ?)")
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
	rows, err := DB.Query(`SELECT username FROM users`)
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
	rows, err := DB.Query(`SELECT email FROM users`)
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
	var user User
	err := DB.QueryRow(`SELECT id, email, username, hash FROM users WHERE email = ?`, email).Scan(&user.Id, &user.Email, &user.Username, &user.Hash)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func getUserBySession(sessionValue string) (User, error) {
	var user User
	err := DB.QueryRow(`SELECT id, email, username, hash FROM users WHERE session_key = ?`, sessionValue).Scan(&user.Id, &user.Email, &user.Username, &user.Hash)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func getUserById(id int) (User, error) {
	var user User
	err := DB.QueryRow(`SELECT username FROM users WHERE id = ?`, id).Scan(&user.Username)
	if err != nil {
		return User{}, err
	}
	user.Id = id
	return user, nil
}

func checkIfUserAlreadyLikedPost(userId, postId int) (int, error) {
	var existingReactionId int
	err := DB.QueryRow(`
		SELECT id FROM reactions
		WHERE user_id = ? AND post_id = ? AND value = 1
		`, userId, postId).Scan(&existingReactionId)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return existingReactionId, nil
}

func checkIfUserAlreadyDislikedPost(userId, postId int) (int, error) {
	var existingDislikeId int
	err := DB.QueryRow(`
		SELECT id FROM reactions
		WHERE user_id = ? AND post_id = ? AND value = 2
		`, userId, postId).Scan(&existingDislikeId)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return existingDislikeId, nil
}

func addLikeToPost(userId, postId int) error {
	_, err := DB.Exec(`
		INSERT INTO reactions (user_id, post_id, value)
		VALUES (?, ?, 1)
		`, userId, postId)
	if err != nil {
		return err
	}
	return nil
}

func removeDislikeFromPost(dislikeId int) error {
	_, err := DB.Exec(`
		DELETE FROM reactions
		WHERE id = ?
		`, dislikeId)
	if err != nil {
		return err
	}
	return nil
}

func addDislikeToPost(userId, postId int) error {
	_, err := DB.Exec(`
		INSERT INTO reactions (user_id, post_id, value)
		VALUES (?, ?, 2)
		`, userId, postId)
	if err != nil {
		return err
	}
	return nil
}

func removeLikeFromPost(userId, postId int) error {
	_, err := DB.Exec(`
		DELETE FROM reactions
		WHERE user_id = ? AND post_id = ? AND value = 1
		`, userId, postId)
	if err != nil {
		return err
	}
	return nil
}

func checkIfUserAlreadyLikedComment(userId, commentId int) (int, error) {
	var existingReactionId int
	err := DB.QueryRow(`
		SELECT id FROM reactions
		WHERE user_id = ? AND comment_id = ? AND value = 1
		`, userId, commentId).Scan(&existingReactionId)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return existingReactionId, nil
}

func checkIfUserAlreadyDislikedComment(userId, commentId int) (int, error) {
	var existingDislikeId int
	err := DB.QueryRow(`
		SELECT id FROM reactions
		WHERE user_id = ? AND comment_id = ? AND value = 2
		`, userId, commentId).Scan(&existingDislikeId)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return existingDislikeId, nil
}

func addLikeToComment(userId, commentId int) error {
	_, err := DB.Exec(`
		INSERT INTO reactions (user_id, comment_id, value)
		VALUES (?, ?, 1)
		`, userId, commentId)
	if err != nil {
		return err
	}
	return nil
}

func addDislikeToComment(userId, commentId int) error {
	_, err := DB.Exec(`
		INSERT INTO reactions (user_id, comment_id, value)
		VALUES (?, ?, 2)
		`, userId, commentId)
	if err != nil {
		return err
	}
	return nil
}

func removeDislikeFromComment(dislikeId int) error {
	_, err := DB.Exec(`
		DELETE FROM reactions
		WHERE id = ?
		`, dislikeId)
	if err != nil {
		return err
	}
	return nil
}

func removeLikeFromComment(userId, commentId int) error {
	_, err := DB.Exec(`
		DELETE FROM reactions
		WHERE user_id = ? AND comment_id = ? AND value = 1
		`, userId, commentId)
	if err != nil {
		return err
	}
	return nil
}

func setUserSession(id int, session_key string) error {
	stmt, err := DB.Prepare("UPDATE users SET session_key = ? WHERE id = ?")
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
	rows, err := DB.Query(`SELECT id, name FROM categories`)
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
	var category Category
	err := DB.QueryRow(`SELECT name FROM categories WHERE id = ?`, id).Scan(&category.Name)
	if err != nil {
		return Category{}, err
	}
	category.Id = id
	return category, nil
}

func getCategoriesByPostId(post_id int) (Categories, error) {
	var categories Categories
	rows, err := DB.Query(`
	SELECT c.id, c.name
	FROM categories c
	JOIN posts_categories pc ON c.id = pc.category_id
	WHERE pc.post_id = ?
	`, post_id)
	if err != nil {
		return Categories{}, err
	}
	defer rows.Close()
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

func addCategory(category Category) error {
	err := category.validateCategory()
	if err != nil {
		return err
	}
	if (&category).DoesCategoryExist() {
		return ErrorCategoryAlreadyExists
	}
	stmt, err := DB.Prepare("INSERT INTO categories (name) VALUES (?)")
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
	var posts Posts
	rows, err := DB.Query(`
	SELECT posts.id, posts.title, posts.body, posts.timestamp
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
		var ts string
		err = rows.Scan(&post.Id, &post.Title, &post.Body, &ts)
		if err != nil {
			return Posts{}, err
		}
		t, err := convertStringToTime(ts)
		if err != nil {
			return Posts{}, err
		}
		post.TimestampString = convertTimeToString(t)
		posts = append(posts, post)
	}
	return posts, nil
}

func getLikedPostsByUserId(user_id int) (Posts, error) {
	var posts Posts
	rows, err := DB.Query(`
	SELECT posts.id, posts.title, posts.body, posts.timestamp
	FROM posts
	JOIN reactions r ON posts.user_id = r.user_id
	WHERE posts.user_id = ? AND r.value = 1
	`, user_id)
	if err != nil {
		return Posts{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var post Post
		var ts string
		err = rows.Scan(&post.Id, &post.Title, &post.Body, &ts)
		if err != nil {
			return Posts{}, err
		}
		t, err := convertStringToTime(ts)
		if err != nil {
			return Posts{}, err
		}
		post.TimestampString = convertTimeToString(t)
		posts = append(posts, post)
	}
	return posts, nil
}

func getPostsByUserId(id int) (Posts, error) {
	var posts Posts
	rows, err := DB.Query(`
	SELECT posts.id, posts.title, posts.body, posts.timestamp
	FROM posts
	WHERE user_id = ?`, id)
	if err != nil {
		return Posts{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var post Post
		var ts string
		err = rows.Scan(&post.Id, &post.Title, &post.Body, &ts)
		if err != nil {
			return Posts{}, err
		}
		t, err := convertStringToTime(ts)
		if err != nil {
			return Posts{}, err
		}
		post.TimestampString = convertTimeToString(t)
		posts = append(posts, post)
	}
	return posts, nil
}

func getAllPosts() (Posts, error) {
	rows, err := DB.Query(`SELECT id, title, body, timestamp FROM posts`)
	if err != nil {
		return Posts{}, err
	}
	defer rows.Close()
	var posts Posts
	for rows.Next() {
		var post Post
		var ts string
		err = rows.Scan(&post.Id, &post.Title, &post.Body, &ts)
		if err != nil {
			return Posts{}, err
		}
		t, err := convertStringToTime(ts)
		if err != nil {
			return Posts{}, err
		}
		post.TimestampString = convertTimeToString(t)
		posts = append(posts, post)
	}
	return posts, nil
}

func getPostById(id int) (Post, error) {
	var post Post
	var ts string
	err := DB.QueryRow(`SELECT title, body, timestamp, user_id FROM posts WHERE id = ?`, id).Scan(&post.Title, &post.Body, &ts, &post.UserId)
	if err != nil {
		return Post{}, err
	}
	t, err := convertStringToTime(ts)
	if err != nil {
		return Post{}, err
	}
	post.TimestampString = convertTimeToString(t)
	post.Id = id
	post.User, err = getUserById(post.UserId)
	if err != nil {
		return Post{}, err
	}
	return post, nil
}

func addPost(post Post) (int, error) {
	err := post.validatePost()
	if err != nil {
		return 0, err
	}
	stmt, err := DB.Prepare("INSERT INTO posts (title, body, user_id, timestamp) VALUES (?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(post.Title, post.Body, post.UserId, getCurrentTimestamp())
	if err != nil {
		return 0, err
	}
	// Get the last inserted post ID
	postId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	for _, category := range post.Categories {
		err = addCategoryToPost(postId, category.Id)
		if err != nil {
			return 0, err
		}
	}
	return int(postId), nil
}

func addCategoryToPost(post_id int64, category_id int) error {
	stmt, err := DB.Prepare("INSERT INTO posts_categories (post_id, category_id) VALUES (?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(post_id, category_id)
	if err != nil {
		return err
	}
	return nil
}

func addComment(comment Comment) (int, error) {
	if err := comment.validateComment(); err != nil {
		return 0, err
	}
	res, err := DB.Exec(
		"INSERT INTO comments (post_id, user_id, body, timestamp) VALUES (?, ?, ?, ?)",
		comment.PostId,
		comment.UserId,
		comment.Body,
		getCurrentTimestamp(),
	)
	commentId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(commentId), nil
}

func getCommentsByPostId(postId int) (Comments, error) {
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
			return Comments{}, err
		}
		t, err := convertStringToTime(ts)
		if err != nil {
			return Comments{}, err
		}
		comment.TimestampString = convertTimeToString(t)
		comments = append(comments, comment)
	}
	return comments, nil
}

func getLikesCountByPostId(postId int) (int, error) {
	var likes int
	err := DB.QueryRow(`
        SELECT COUNT(*)
        FROM reactions
        WHERE post_id = ? AND value = 1
    `, postId).Scan(&likes)
	if err != nil {
		return 0, err
	}
	return likes, nil
}

func getDisikesCountByPostId(postId int) (int, error) {
	var dislikes int
	err := DB.QueryRow(`
        SELECT COUNT(*)
        FROM reactions
        WHERE post_id = ? AND value = 2
    `, postId).Scan(&dislikes)
	if err != nil {
		return 0, err
	}
	return dislikes, nil
}

func getLikesCountByCommentId(commentId int) (int, error) {
	var likes int
	err := DB.QueryRow(`
        SELECT COUNT(*)
        FROM reactions
        WHERE comment_id = ? AND value = 1
    `, commentId).Scan(&likes)
	if err != nil {
		return 0, err
	}
	return likes, nil
}

func getDisikesCountByCommentId(commentId int) (int, error) {
	var dislikes int
	err := DB.QueryRow(`
        SELECT COUNT(*)
        FROM reactions
        WHERE comment_id = ? AND value = 2
    `, commentId).Scan(&dislikes)
	if err != nil {
		return 0, err
	}
	return dislikes, nil
}
