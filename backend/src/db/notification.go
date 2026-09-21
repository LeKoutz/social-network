package db

import (
	"forum/src/ferror"
)

type NotificationRowType struct {
	Id              int64
	UserId          int64
	ActorId         int64
	Type            string
	PostId          int64
	CommentId       int64
	Timestamp       string
	Read            bool
	Username        string
}

func (n *NotificationRowType) InsertNotification() error {
	var commentId any
	if n.CommentId == 0 {
		commentId = nil
	} else {
		commentId = n.CommentId
	}
	query := `INSERT INTO notifications (user_id, actor_id, type, post_id, comment_id, timestamp) VALUES (?, ?, ?, ?, ?, ?)`
	res, err := db.Exec(query, n.UserId, n.ActorId, n.Type, n.PostId, commentId, n.Timestamp)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	n.Id, err = res.LastInsertId()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}

func (u *UserRowType) UpdateNotificationAsRead(notificationId int64) error {
	stmt, err := db.Prepare(`UPDATE notifications SET "read" = 1 WHERE id = ? AND user_id = ?`)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	_, err = stmt.Exec(notificationId, u.Id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}
