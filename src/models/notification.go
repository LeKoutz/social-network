package models

import (
	"errors"
	"forum/src/db"
	"forum/src/utils"
)

type NotificationType struct {
	db.NotificationRowType
	Actor   UserType
	Post    PostType
	Comment CommentType
}

func CreateNotification(notification NotificationType) error {
	if notification.ActorId == notification.UserId {
		return nil
	}
	notification.TimestampString = utils.GetCurrentTimestamp()
	err := notification.InsertNotification()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}

func (user *UserType) MarkAsReadPost(post PostType) error {
	for i, notification := range user.Notifications {
		if notification.PostId == post.Id && !notification.Read {
			err := user.UpdateNotificationAsRead(notification.Id)
			if err != nil {
				if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
				return err
			}
			user.Notifications[i].Read = true
			user.UnreadNotificationsCount--
		}
	}
	return nil
}

func (comment *CommentType) CreateCommentNotification(post PostType) error {
	var notification NotificationType
	notification.UserId = post.UserId
	notification.ActorId = comment.UserId
	notification.Type = "comment"
	notification.PostId = comment.PostId
	notification.CommentId = comment.Id
	return CreateNotification(notification)
}

func (comment *CommentType) CreateReactionNotification(userId int64, t string) error {
	var notification NotificationType
	notification.UserId = comment.UserId
	notification.ActorId = userId
	notification.CommentId = comment.Id
	notification.PostId = comment.PostId
	notification.Type = t
	return CreateNotification(notification)
}

func (post *PostType) CreateReactionNotification(userId int64, t string) error {
	var notification NotificationType
	notification.UserId = post.UserId
	notification.ActorId = userId
	notification.Type = t
	notification.PostId = post.Id
	return CreateNotification(notification)
}
