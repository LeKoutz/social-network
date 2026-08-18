package models

import (
	"errors"
	"forum/src/db"
	"forum/src/ferror"
	"forum/src/utils"
	"net/mail"
	"regexp"
	"slices"
	"sort"
)

type UserType struct {
	db.UserRowType

	LoggedIn                 bool
	Notifications            NotificationsType
	UnreadNotificationsCount int
	Activities               ActivitiesType
	ChatMessages             ChatMessagesType
	// LastMessageTimestamp     int64
	Identifier               string
	Password				 string
}

func GetGuestUser() UserType {
	var u UserType
	u.Username = "guest"
	u.LoggedIn = false
	return u
}

func (u *UserType) ValidateUsername() error {
	unameMask := regexp.MustCompile(`^[a-zA-Z0-9_]{4,50}$`)
	if !unameMask.MatchString((*u).Username) {
		return ferror.ErrorInvalidUsername
	}
	return nil
}

func (u *UserType) ValidateEmail() error {
	_, err := mail.ParseAddress(u.Email)
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	return nil
}

func (u *UserType) ValidateUser() error {
	var err error
	if err = u.ValidateUsername(); err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	if !IsUniqueUsername(u.Username) {
		return ferror.ErrorUsernameTaken
	}
	if err = u.ValidateEmail(); err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	if IsEmailRegistered(u.Email) {
		return ferror.ErrorEmailIsRegistered
	}
	if !u.isValidGender() {
		return ferror.ErrorInvalidGender
	}
	return nil
}

func (u *UserType) isValidGender() bool {
	if u.Gender != "male" && u.Gender != "female" && u.Gender != "other" {
		return false
	}
	return true
}

func (u *UserType) Add() error {
	err := u.ValidateUser()
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	return u.InsertUserWithHash()
}

func (u *UserType) AddOAuth() error {
	if err := u.ValidateUser(); err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	if IsEmailRegistered(u.Email) {
		return ferror.ErrorEmailIsRegistered
	}
	if !IsUniqueUsername(u.Username) {
		return ferror.ErrorUsernameTaken
	}
	return u.InsertUserWithOAuth()
}

func (u *UserType) GetPosts() (PostsType, error) {
	var posts PostsType
	rows, err := u.SelectPosts()
	if err != nil {
		if config.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return posts, err
	}
	for _, row := range rows {
		var post PostType
		post.PostRowType = row
		posts = append(posts, post)
	}
	return posts, err
}

func (u *UserType) GetLikedPosts() (PostsType, error) {
	var posts PostsType
	rows, err := u.SelectLikedPosts()
	if err != nil {
		if config.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return posts, err
	}
	for _, row := range rows {
		var post PostType
		post.PostRowType = row
		posts = append(posts, post)
	}
	return posts, err
}

func (u *UserType) GetUserBySession() error {
	err := u.SelectUserBySession()
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	return nil
}

func (u *UserType) GetUserByIdentifier() error {
	err := u.SelectUserByIdentifier(u.Identifier)
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return err
	}
	return nil
}

func (u *UserType) GetUserPasswordByIdentifier() error {
	err := u.SelectUserPasswordByIdentifier(u.Identifier)
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return err
	}
	return nil
}

func (u *UserType) GetUserByOAuthProviderAndEmail() error {
	err := u.SelectUserByOAuthProviderAndEmail()
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return err
	}
	return nil
}

func (u *UserType) SetUserSession(session_key string) error {
	var err error
	err = u.UpdateUserSession(session_key)
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return err
	}
	u.SessionId = session_key
	return nil
}

func IsUniqueUsername(username string) bool {
	usernames, err := db.SelectAllUsernames()
	if err != nil {
		(&ferror.Error{}).Consume(err).LogError()
		return false
	}
	return !slices.Contains(usernames, username)
}

func IsUniqueEmail(email string) bool {
	emails, err := db.SelectAllUserEmails()
	if err != nil {
		(&ferror.Error{}).Consume(err).LogError()
		return false
	}
	return !slices.Contains(emails, email)
}

func IsEmailRegistered(email string) bool {
	return !IsUniqueEmail(email)
}

func IsUsernameRegistered(username string) bool {
	return !IsUniqueUsername(username)
}

// Check if user already liked this post
func HasUserLikedPost(userId, postId int64) (bool, error) {
	reactionId, err := db.SelectUserLikeFromPost(userId, postId)
	if err != nil {
		return false, err
	}
	return reactionId != 0, nil
}

// Check if user already disliked this post
func HasUserDislikedPost(userId, postId int64) (bool, error) {
	reactionId, err := db.SelectUserDislikeFromPost(userId, postId)
	if err != nil {
		return false, err
	}
	return reactionId != 0, nil
}

// Check if user already liked this comment
func HasUserLikedComment(userId, commentId int64) (bool, error) {
	reactionId, err := db.SelectUserLikeFromComment(userId, commentId)
	if err != nil {
		return false, err
	}
	return reactionId != 0, nil
}

// Check if user already disliked this comment
func HasUserDislikedComment(userId, commentId int64) (bool, error) {
	reactionId, err := db.SelectUserDislikeFromComment(userId, commentId)
	if err != nil {
		return false, err
	}
	return reactionId != 0, nil
}

func (u *UserType) GetNotifications() error {
	rows, err := db.SelectNotificationsByUserId(u.Id)
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return err
	}
	for _, row := range rows {
		var notification NotificationType
		notification.NotificationRowType = row
		u.Notifications = append(u.Notifications, notification)
	}
	u.CountUnreadNotifications()
	return nil
}

func (u *UserType) MarkNotificationAsRead(notificationId int64) error {
	err := u.UpdateNotificationAsRead(notificationId)
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return err
	}
	return nil
}

func (u *UserType) MarkAllNotificationsAsRead() error {
	err := u.UpdateAllNotificationsAsRead()
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return err
	}
	return nil
}

func (u *UserType) GetActivity() error {
	err := u.GetPostsActivity()
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return err
	}
	err = u.GetCommentsActivity()
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return err
	}
	err = u.GetLikedPostsActivity()
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return err
	}
	err = u.GetDislikedPostsActivity()
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return err
	}
	err = u.GetLikedCommentsActivity()
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return err
	}
	err = u.GetDislikedCommentsActivity()
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return err
	}
	sort.Slice(u.Activities, func(i, j int) bool {
		return u.Activities[i].TimestampString > u.Activities[j].TimestampString
	})
	return nil
}

func (u *UserType) GetPostsActivity() error {
	posts, err := u.GetPosts()
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return err
	}
	for _, post := range posts {
		err := post.SelectPostById()
		if err != nil {
			if config.Debug {
				err = errors.Join(utils.GetFunctionName(), err)
			}
			return err
		}
		err = post.GetReactions()
		if err != nil {
			if config.Debug {
				err = errors.Join(utils.GetFunctionName(), err)
			}
			return err
		}
		err = post.GetReactionsByUserId(u.Id)
		if err != nil {
			if config.Debug {
				err = errors.Join(utils.GetFunctionName(), err)
			}
			return err
		}
		var activity ActivityType
		activity.TimestampString = post.TimestampString
		activity.Post = post
		activity.Type = "post"
		u.Activities = append(u.Activities, activity)
	}
	return nil
}

func (u *UserType) GetCommentsActivity() error {
	rows, err := db.SelectCommentsByUserId(u.Id)
	if err != nil {
		if config.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	for _, row := range rows {
		var activity ActivityType
		var post PostType
		var comment CommentType
		comment.CommentRowType = row
		post.Id = comment.PostId

		err = post.SelectPostById()
		if err != nil {
			if config.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return err
		}
		err = comment.GetReactions()
		if err != nil {
			if config.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return err
		}
		err = comment.GetReactionsByUserId(u.Id)
		if err != nil {
			if config.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return err
		}
		activity.Type = "comment"
		activity.Comment = comment
		activity.TimestampString = comment.TimestampString
		activity.Post = post
		u.Activities = append(u.Activities, activity)
	}
	return nil
}

func (u *UserType) GetLikedPostsActivity() error {
	var reactions ReactionsType
	err := reactions.GetPostLikesByUserId(u.Id)
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return err
	}
	for _, reaction := range reactions {
		var activity ActivityType
		var post PostType
		post.Id = reaction.PostId
		err = post.SelectPostById()
		if err != nil {
			if config.Debug {
				err = errors.Join(utils.GetFunctionName(), err)
			}
			return err
		}
		err = post.GetReactions()
		if err != nil {
			if config.Debug {
				err = errors.Join(utils.GetFunctionName(), err)
			}
			return err
		}
		err = post.GetReactionsByUserId(u.Id)
		if err != nil {
			if config.Debug {
				err = errors.Join(utils.GetFunctionName(), err)
			}
			return err
		}
		activity.Type = "postLike"
		activity.TimestampString = reaction.TimestampString
		activity.Post = post
		u.Activities = append(u.Activities, activity)
	}
	return nil
}

func (u *UserType) GetDislikedPostsActivity() error {
	reactions, err := GetPostDislikesByUserId(u.Id)
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return err
	}
	for _, reaction := range reactions {
		var activity ActivityType
		var post PostType
		post.Id = reaction.PostId
		err = post.SelectPostById()
		if err != nil {
			if config.Debug {
				err = errors.Join(utils.GetFunctionName(), err)
			}
			return err
		}
		err = post.GetReactions()
		if err != nil {
			if config.Debug {
				err = errors.Join(utils.GetFunctionName(), err)
			}
			return err
		}
		err = post.GetReactionsByUserId((*u).Id)
		if err != nil {
			if config.Debug {
				err = errors.Join(utils.GetFunctionName(), err)
			}
			return err
		}
		activity.Type = "postDislike"
		activity.TimestampString = reaction.TimestampString
		activity.Post = post
		u.Activities = append(u.Activities, activity)
	}
	return nil
}

func (u *UserType) GetLikedCommentsActivity() error {
	reactions, err := GetCommentLikesByUserId(u.Id)
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return err
	}
	for _, reaction := range reactions {
		var activity ActivityType
		var comment CommentType
		comment.Id = reaction.CommentId
		err = comment.SelectCommentById()
		if err != nil {
			if config.Debug {
				err = errors.Join(utils.GetFunctionName(), err)
			}
			return err
		}
		err = comment.GetReactions()
		if err != nil {
			if config.Debug {
				err = errors.Join(utils.GetFunctionName(), err)
			}
			return err
		}
		err = comment.GetReactionsByUserId(u.Id)
		if err != nil {
			if config.Debug {
				err = errors.Join(utils.GetFunctionName(), err)
			}
			return err
		}
		var post PostType
		post.Id = comment.PostId
		err = post.SelectPostById()
		if err != nil {
			if config.Debug {
				err = errors.Join(utils.GetFunctionName(), err)
			}
			return err
		}
		activity.Type = "commentLike"
		activity.TimestampString = reaction.TimestampString
		activity.Comment = comment
		activity.Post = post
		u.Activities = append(u.Activities, activity)
	}
	return nil
}

func (u *UserType) GetDislikedCommentsActivity() error {
	reactions, err := GetCommentDisikesByUserId(u.Id)
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return err
	}
	for _, reaction := range reactions {
		var activity ActivityType
		var comment CommentType
		comment.Id = reaction.CommentId
		err = comment.SelectCommentById()
		if err != nil {
			if config.Debug {
				err = errors.Join(utils.GetFunctionName(), err)
			}
			return err
		}
		err = comment.GetReactions()
		if err != nil {
			if config.Debug {
				err = errors.Join(utils.GetFunctionName(), err)
			}
			return err
		}
		err = comment.GetReactionsByUserId(u.Id)
		if err != nil {
			if config.Debug {
				err = errors.Join(utils.GetFunctionName(), err)
			}
			return err
		}
		var post PostType
		post.Id = comment.PostId
		err = post.SelectPostById()
		if err != nil {
			if config.Debug {
				err = errors.Join(utils.GetFunctionName(), err)
			}
			return err
		}
		activity.Type = "commentDislike"
		activity.TimestampString = reaction.TimestampString
		activity.Comment = comment
		activity.Post = post
		u.Activities = append(u.Activities, activity)
	}
	return nil
}

func (u *UserType) CountUnreadNotifications() {
	for _, notification := range u.Notifications {
		if !notification.Read {
			u.UnreadNotificationsCount++
		}
	}
}
