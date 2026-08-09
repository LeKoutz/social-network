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

// func (n *NotificationType) ToNotificationRowType() *db.NotificationRowType {
// 	return &db.NotificationRowType{
// 		UserId:          n.UserId,
// 		ActorId:         n.ActorId,
// 		Type:            n.Type,
// 		PostId:          n.PostId,
// 		CommentId:       n.CommentId,
// 		TimestampString: n.TimestampString,
// 	}
// }

func (n *NotificationType) FromNotificationRowType(nr *db.NotificationRowType) {
	n.Id = nr.Id
	n.UserId = nr.UserId
	n.ActorId = nr.ActorId
	n.Type = nr.Type
	n.PostId = nr.PostId
	n.CommentId = nr.CommentId
	n.Read = nr.Read
	var user UserType
	user.Username = nr.Username
	n.Actor = user
	n.TimestampString = nr.TimestampString
}

func CreateNotification(notification NotificationType) error {
	if notification.ActorId == notification.UserId {
		return nil
	}
	notification.TimestampString = utils.GetCurrentTimestamp()
	return notification.InsertNotification()
}

func (user *UserType) MarkAsReadPost(post PostType) error {
	for i, notification := range user.Notifications {
		if notification.PostId == post.Id && !notification.Read {
			err := user.UpdateNotificationAsRead(notification.Id)
			if err != nil {
				return errors.Join(utils.GetFunctionName(), err)
			}
			user.Notifications[i].Read = true
			user.UnreadNotificationsCount--
		}
	}
	return nil
}

func (comment *CommentType) CreateCommentNotification(post PostType) error {
	var notification NotificationType
	notification.UserId = post.User.Id
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
	notification.UserId = post.User.Id
	notification.ActorId = userId
	notification.Type = t
	notification.PostId = post.Id
	return CreateNotification(notification)
}
