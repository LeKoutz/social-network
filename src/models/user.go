package models

import (
	"errors"
	"forum/src/utils"
	"net/mail"
	"regexp"
	"slices"
)

type User struct {
	Id            int64
	Username      string
	Hash          string
	Email         string
	Role          string
	LoggedIn      bool
	OAuthProvider string
	OAuthId		  string
	OwnedPosts    Posts
	OwnedComments Comments
	OwnedLikes    Likes
	OwnedDislikes Dislikes
}

func GetGuestUser() User {
	return User{
		Username: "guest",
		Role:     "guest",
		LoggedIn: false,
	}
}

// Returns the user, if any, for a given email address
// Populates the `User.Hash` field for comparison against the given password
func GetUserByEmail(email string) (User, error) {
	var user User
	err := DB.QueryRow(`SELECT id, email, username, hash FROM users WHERE email = ?`, email).Scan(&user.Id, &user.Email, &user.Username, &user.Hash)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return User{}, err
	}
	return user, nil
}

// Returns the user's id, email and username for given session ID
func GetUserBySession(sessionValue string) (User, error) {
	var user User
	err := DB.QueryRow(`SELECT id, email, username FROM users WHERE session_key = ?`, sessionValue).Scan(&user.Id, &user.Email, &user.Username)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return User{}, err
	}
	return user, nil
}

// WIP
// Returns a user given the oauthID and provider
func GetUserByOAuth(provider, oauthID string) (User, error) {
    var user User
	err := DB.QueryRow(`SELECT id, email, username, oauth_provider, oauth_id FROM users WHERE oauth_provider = ? AND oauth_id = ?`, provider, oauthID).Scan(&user.Id, &user.Email, &user.Username, &user.OAuthProvider, &user.OAuthId)
    if err != nil {
        return User{}, err
    }
    return user, nil
}

func getUserById(id int64) (User, error) {
	var user User
	err := DB.QueryRow(`SELECT username FROM users WHERE id = ?`, id).Scan(&user.Username)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return User{}, err
	}
	user.Id = id
	return user, nil
}

func (u *User) ValidateUsername() error {
	unameMask := regexp.MustCompile(`^[a-zA-Z0-9]{4,15}$`)
	if !unameMask.MatchString((*u).Username) {
		return ErrorInvalidUsername
	}
	return nil
}

func (u *User) ValidateEmail() error {
	_, err := mail.ParseAddress(u.Email)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) ValidateUser() error {
	var err error
	if err = u.ValidateUsername(); err != nil {
		return err
	}
	if err = u.ValidateEmail(); err != nil {
		return err
	}
	return nil
}

func (u *User) Add() error {
	err := u.ValidateUser()
	if err != nil {
		return err
	}
	if IsEmailRegistered(u.Email) {
		return ErrorEmailIsRegistered
	}
	stmt, err := DB.Prepare("INSERT INTO users (username, email, hash) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(u.Username, u.Email, u.Hash)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) AddOAuth() error {
    if err := u.ValidateUser(); 
		err != nil {
        return err
    }
    if IsEmailRegistered(u.Email) {
        return ErrorEmailIsRegistered
    }
    if !IsUniqueUsername(u.Username) {
        return ErrorUsernameTaken
    }
	stmt, err := DB.Prepare("INSERT INTO users (username, email, oauth_provider, oauth_id) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
    _, err = stmt.Exec(u.Username, u.Email, u.OAuthProvider, u.OAuthId)
	if err != nil {
        return err
    }

    return nil
}

func (u *User) GetPosts() (Posts, error) {
	var posts Posts
	rows, err := DB.Query(`
	SELECT posts.id, posts.title, posts.body, posts.timestamp
	FROM posts
	WHERE user_id = ?`, (*u).Id)
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

func (u *User) GetLikedPosts() (Posts, error) {
	var posts Posts
	rows, err := DB.Query(`
	SELECT posts.id, posts.title, posts.body, posts.timestamp
	FROM posts
	JOIN reactions r ON posts.id = r.post_id
	WHERE r.user_id = ? AND r.value = 1
	`, (*u).Id)
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

func (u *User) SetUserSession(session_key string) error {
	stmt, err := DB.Prepare("UPDATE users SET session_key = ? WHERE id = ?")
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	_, err = stmt.Exec(session_key, (*u).Id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	return nil
}

func IsUniqueUsername(username string) bool {
	usernames, err := GetAllUsernames()
	if err != nil {
		(&Error{}).Consume(err).LogError()
		return false
	}
	return !slices.Contains(usernames, username)
}

func IsUniqueEmail(email string) bool {
	emails, err := GetAllUserEmails()
	if err != nil {
		(&Error{}).Consume(err).LogError()
		return false
	}
	return !slices.Contains(emails, email)
}

func IsEmailRegistered(email string) bool {
	return !IsUniqueEmail(email)
}

// Check if user already liked this post
func HasUserLikedPost(userId, postId int64) (bool, error) {
	reactionId, err := CheckIfUserLikedPost(userId, postId)
	if err != nil {
		return false, err
	}
	return reactionId != 0, nil
}

// Check if user already disliked this post
func HasUserDislikedPost(userId, postId int64) (bool, error) {
	reactionId, err := CheckIfUserDislikedPost(userId, postId)
	if err != nil {
		return false, err
	}
	return reactionId != 0, nil
}

// Check if user already liked this comment
func HasUserLikedComment(userId, commentId int64) (bool, error) {
	reactionId, err := CheckIfUserLikedComment(userId, commentId)
	if err != nil {
		return false, err
	}
	return reactionId != 0, nil
}

// Check if user already disliked this comment
func HasUserDislikedComment(userId, commentId int64) (bool, error) {
	reactionId, err := CheckIfUserDislikedComment(userId, commentId)
	if err != nil {
		return false, err
	}
	return reactionId != 0, nil
}
