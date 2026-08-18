package db

import (
	"errors"
	"forum/src/utils"
)

type NotificationRowType struct {
	Id              int64
	UserId          int64
	ActorId         int64
	Type            string
	PostId          int64
	CommentId       int64
	TimestampString string
	Read            bool
	Username        string
}

func (n *NotificationRowType) InsertNotification() error {
	query := `INSERT INTO notifications (user_id, actor_id, type, post_id, comment_id, timestamp) VALUES (?, ?, ?, ?, ?, ?)`
	res, err := db.Exec(query, n.UserId, n.ActorId, n.Type, n.PostId, n.CommentId, n.TimestampString)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	n.Id, err = res.LastInsertId()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}

func (u *UserRowType) UpdateNotificationAsRead(notificationId int64) error {
	stmt, err := db.Prepare(`UPDATE notifications SET "read" = 1 WHERE id = ? AND user_id = ?`)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	_, err = stmt.Exec(notificationId, u.Id)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}
