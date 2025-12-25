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

func GetUserByEmail(email string) (User, error) {
	var user User
	err := DB.QueryRow(`SELECT id, email, username, hash FROM users WHERE email = ?`, email).Scan(&user.Id, &user.Email, &user.Username, &user.Hash)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return User{}, err
	}
	return user, nil
}

func GetUserBySession(sessionValue string) (User, error) {
	var user User
	err := DB.QueryRow(`SELECT id, email, username, hash FROM users WHERE session_key = ?`, sessionValue).Scan(&user.Id, &user.Email, &user.Username, &user.Hash)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
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

func RegisterUserOnDB(user User) error {
	err := user.ValidateUser()
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

func IsUniqueUsername(username string) bool {
	usernames, err := GetAllUsernames()
	if err != nil {
		var e Error
		e.Consume(err)
		e.LogError()
		return false
	}
	return !slices.Contains(usernames, username)
}

func IsUniqueEmail(email string) bool {
	emails, err := GetAllUserEmails()
	if err != nil {
		var e Error
		e.Consume(err)
		e.LogError()
		return false
	}
	return !slices.Contains(emails, email)
}

func IsEmailRegistered(email string) bool {
	return !IsUniqueEmail(email)
}

// Check if user already liked this post
func HasUserLikedPost(userId, postId int64) (bool, error) {
	existingReactionId, err := CheckIfUserLikedPost(userId, postId)
	if err != nil {
		return false, err
	}
	return existingReactionId != 0, nil
}

// Check if user already disliked this post
func HasUserDislikedPost(userId, postId int64) (bool, error) {
	existingReactionId, err := CheckIfUserDislikedPost(userId, postId)
	if err != nil {
		return false, err
	}
	return existingReactionId != 0, nil
}

// Check if user already liked this comment
func HasUserLikedComment(userId, commentId int64) (bool, error) {
	existingReactionId, err := CheckIfUserLikedComment(userId, commentId)
	if err != nil {
		return false, err
	}
	return existingReactionId != 0, nil
}

// Check if user already disliked this comment
func HasUserDislikedComment(userId, commentId int64) (bool, error) {
	existingReactionId, err := CheckIfUserDislikedComment(userId, commentId)
	if err != nil {
		return false, err
	}
	return existingReactionId != 0, nil
}
