package db

import (
	"errors"
	"forum/src/utils"
)

type NotificationRowsType []NotificationRowType

func GetNotificationsByUserId(userId int64) (NotificationRowsType, error) {
	var notifications NotificationRowsType
	rows, err := db.Query(`
	SELECT n.id, n.user_id, n.actor_id, n.type, n.post_id, comment_id, n.timestamp, n."read", u.username
	FROM notifications n
	JOIN users u ON u.id = n.actor_id
	WHERE user_id = ?
	ORDER BY n.timestamp DESC
	`, userId)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return NotificationRowsType{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var notification NotificationRowType
		err = rows.Scan(&notification.Id,
			&notification.UserId,
			&notification.ActorId,
			&notification.Type,
			&notification.PostId,
			&notification.CommentId,
			&notification.TimestampString,
			&notification.Read,
			&notification.Username)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return NotificationRowsType{}, err
		}
		notifications = append(notifications, notification)
	}
	return notifications, nil
}

func (u *UserRowType) UpdateAllNotificationsAsRead() error {
	stmt, err := db.Prepare(`UPDATE notifications SET "read" = 1 WHERE user_id = ?`)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return errors.Join(utils.GetFunctionName(), err)
	}
	_, err = stmt.Exec((*u).Id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return errors.Join(utils.GetFunctionName(), err)
	}
	return nil
}
