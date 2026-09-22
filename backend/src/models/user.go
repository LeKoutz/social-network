package models

import (
	"forum/src/db"
	"forum/src/ferror"
	"net/mail"
	"regexp"
	"slices"
	"sort"
)

type UserType struct {
	db.UserRowType
	db.UserProfileRowType

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
		err := ferror.ErrorInvalidUsername
		return ferror.ReturnErr(err)
	}
	return nil
}

func (u *UserType) ValidateEmail() error {
	_, err := mail.ParseAddress(u.Email)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}

func (u *UserType) ValidateUser() error {
	var err error
	if err = u.ValidateUsername(); err != nil {
		return ferror.ReturnErr(err)
	}
	if !IsUniqueUsername(u.Username) {
		err = ferror.ErrorUsernameTaken
		return ferror.ReturnErr(err)
	}
	if err = u.ValidateEmail(); err != nil {
		return ferror.ReturnErr(err)
	}
	if IsEmailRegistered(u.Email) {
		err = ferror.ErrorEmailIsRegistered
		return ferror.ReturnErr(err)
	}
	return nil
}

func (u *UserType) Add() error {
	var err error
	err = u.ValidateUser()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	err = u.InsertUserWithHash()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	u.UserId = u.Id
	err = u.InsertUserProfile()
	if err != nil {
		u.DeleteUserById()
		return ferror.ReturnErr(err)
	}
	return nil
}

func (u *UserType) AddOAuth() error {
	var err error
	if err = u.ValidateUser(); err != nil {
		return ferror.ReturnErr(err)
	}
	if IsEmailRegistered(u.Email) {
		err = ferror.ErrorEmailIsRegistered
		return ferror.ReturnErr(err)
	}
	if !IsUniqueUsername(u.Username) {
		err = ferror.ErrorUsernameTaken
		return ferror.ReturnErr(err)
	}
	err = u.InsertUserWithOAuth()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}

func (u *UserType) GetPosts() (PostsType, error) {
	var posts PostsType
	rows, err := u.SelectPosts()
	if err != nil {
		return posts, ferror.ReturnErr(err)
	}
	err = u.GetById()
	if err != nil {
		return posts, ferror.ReturnErr(err)
	}
	for _, row := range rows {
		var post PostType
		post.PostRowType = row
		post.User = *u
		posts = append(posts, post)
	}
	return posts, err
}

func (u *UserType) GetLikedPosts() (PostsType, error) {
	var posts PostsType
	rows, err := u.SelectLikedPosts()
	if err != nil {
		return posts, ferror.ReturnErr(err)
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
		return ferror.ReturnErr(err)
	}
	return nil
}

func (u *UserType) GetUserByIdentifier() error {
	err := u.SelectUserByIdentifier(u.Identifier)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}

func (u *UserType) GetUserPasswordByIdentifier() error {
	err := u.SelectUserPasswordByIdentifier(u.Identifier)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}

func (u *UserType) GetUserByOAuthProviderAndEmail() error {
	err := u.SelectUserByOAuthProviderAndEmail()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}

func (u *UserType) SetUserSession(session_key string) error {
	var err error
	err = u.UpdateUserSession(session_key)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	u.SessionId = session_key
	return nil
}

func IsUniqueUsername(username string) bool {
	usernames, err := db.SelectAllUsernames()
	if err != nil {
		(&ferror.Error{}).Consume(ferror.ReturnErr(err)).LogError()
		return false
	}
	return !slices.Contains(usernames, username)
}

func IsUniqueEmail(email string) bool {
	emails, err := db.SelectAllUserEmails()
	if err != nil {
		(&ferror.Error{}).Consume(ferror.ReturnErr(err)).LogError()
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
		return false, ferror.ReturnErr(err)
	}
	return reactionId != 0, nil
}

// Check if user already disliked this post
func HasUserDislikedPost(userId, postId int64) (bool, error) {
	reactionId, err := db.SelectUserDislikeFromPost(userId, postId)
	if err != nil {
		return false, ferror.ReturnErr(err)
	}
	return reactionId != 0, nil
}

// Check if user already liked this comment
func HasUserLikedComment(userId, commentId int64) (bool, error) {
	reactionId, err := db.SelectUserLikeFromComment(userId, commentId)
	if err != nil {
		return false, ferror.ReturnErr(err)
	}
	return reactionId != 0, nil
}

// Check if user already disliked this comment
func HasUserDislikedComment(userId, commentId int64) (bool, error) {
	reactionId, err := db.SelectUserDislikeFromComment(userId, commentId)
	if err != nil {
		return false, ferror.ReturnErr(err)
	}
	return reactionId != 0, nil
}

func (u *UserType) GetNotifications() error {
	rows, err := db.SelectNotificationsByUserId(u.Id)
	if err != nil {
		return ferror.ReturnErr(err)
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
		return ferror.ReturnErr(err)
	}
	return nil
}

func (u *UserType) MarkAllNotificationsAsRead() error {
	err := u.UpdateAllNotificationsAsRead()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}

func (u *UserType) GetActivity() error {
	err := u.GetPostsActivity()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	err = u.GetCommentsActivity()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	err = u.GetLikedPostsActivity()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	err = u.GetDislikedPostsActivity()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	err = u.GetLikedCommentsActivity()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	err = u.GetDislikedCommentsActivity()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	sort.Slice(u.Activities, func(i, j int) bool {
		return u.Activities[i].Timestamp > u.Activities[j].Timestamp
	})
	return nil
}

func (u *UserType) GetPostsActivity() error {
	posts, err := u.GetPosts()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	for _, post := range posts {
		err := post.SelectPostById()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		err = post.GetReactions()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		err = post.GetReactionsByUserId(u.Id)
		if err != nil {
			return ferror.ReturnErr(err)
		}
		var activity ActivityType
		activity.Timestamp = post.Timestamp
		activity.Post = post
		activity.Type = "post"
		u.Activities = append(u.Activities, activity)
	}
	return nil
}

func (u *UserType) GetCommentsActivity() error {
	rows, err := db.SelectCommentsByUserId(u.Id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	for _, row := range rows {
		var activity ActivityType
		var post PostType
		var comment CommentType
		comment.CommentRowType = row
		post.Id = comment.PostId

		err = post.SelectPostById()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		err = comment.GetReactions()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		err = comment.GetReactionsByUserId(u.Id)
		if err != nil {
			return ferror.ReturnErr(err)
		}
		activity.Type = "comment"
		activity.Comment = comment
		activity.Timestamp = comment.Timestamp
		activity.Post = post
		u.Activities = append(u.Activities, activity)
	}
	return nil
}

func (u *UserType) GetLikedPostsActivity() error {
	var reactions ReactionsType
	err := reactions.GetPostLikesByUserId(u.Id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	for _, reaction := range reactions {
		var activity ActivityType
		var post PostType
		post.Id = reaction.PostId
		err = post.SelectPostById()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		err = post.GetReactions()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		err = post.GetReactionsByUserId(u.Id)
		if err != nil {
			return ferror.ReturnErr(err)
		}
		activity.Type = "postLike"
		activity.Timestamp = reaction.Timestamp
		activity.Post = post
		u.Activities = append(u.Activities, activity)
	}
	return nil
}

func (u *UserType) GetDislikedPostsActivity() error {
	reactions, err := GetPostDislikesByUserId(u.Id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	for _, reaction := range reactions {
		var activity ActivityType
		var post PostType
		post.Id = reaction.PostId
		err = post.SelectPostById()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		err = post.GetReactions()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		err = post.GetReactionsByUserId((*u).Id)
		if err != nil {
			return ferror.ReturnErr(err)
		}
		activity.Type = "postDislike"
		activity.Timestamp = reaction.Timestamp
		activity.Post = post
		u.Activities = append(u.Activities, activity)
	}
	return nil
}

func (u *UserType) GetLikedCommentsActivity() error {
	reactions, err := GetCommentLikesByUserId(u.Id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	for _, reaction := range reactions {
		var activity ActivityType
		var comment CommentType
		comment.Id = reaction.CommentId
		err = comment.SelectCommentById()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		err = comment.GetReactions()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		err = comment.GetReactionsByUserId(u.Id)
		if err != nil {
			return ferror.ReturnErr(err)
		}
		var post PostType
		post.Id = comment.PostId
		err = post.SelectPostById()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		activity.Type = "commentLike"
		activity.Timestamp = reaction.Timestamp
		activity.Comment = comment
		activity.Post = post
		u.Activities = append(u.Activities, activity)
	}
	return nil
}

func (u *UserType) GetDislikedCommentsActivity() error {
	reactions, err := GetCommentDisikesByUserId(u.Id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	for _, reaction := range reactions {
		var activity ActivityType
		var comment CommentType
		comment.Id = reaction.CommentId
		err = comment.SelectCommentById()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		err = comment.GetReactions()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		err = comment.GetReactionsByUserId(u.Id)
		if err != nil {
			return ferror.ReturnErr(err)
		}
		var post PostType
		post.Id = comment.PostId
		err = post.SelectPostById()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		activity.Type = "commentDislike"
		activity.Timestamp = reaction.Timestamp
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

func (u *UserType) GetById() error {
	var err error
	err = u.SelectUserById()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}
