package db

import (
	"database/sql"
	"forum/src/ferror"
)

type NotificationRowsType []NotificationRowType

func SelectNotificationsByUserId(userId int64) (NotificationRowsType, error) {
	var notifications NotificationRowsType
	rows, err := db.Query(`
	SELECT n.id, n.user_id, n.actor_id, n.type, n.post_id, comment_id, n.timestamp, n."read", u.username
	FROM notifications n
	JOIN users u ON u.id = n.actor_id
	WHERE user_id = ?
	ORDER BY n.timestamp DESC
	`, userId)
	if err != nil {
		return NotificationRowsType{}, ferror.ReturnErr(err)
	}
	defer rows.Close()
	for rows.Next() {
		var notification NotificationRowType
		var commentId sql.NullInt64
		err = rows.Scan(&notification.Id,
			&notification.UserId,
			&notification.ActorId,
			&notification.Type,
			&notification.PostId,
			&commentId,
			&notification.Timestamp,
			&notification.Read,
			&notification.Username)
		if err != nil {
			return NotificationRowsType{}, ferror.ReturnErr(err)
		}
		notification.CommentId = commentId.Int64
		notifications = append(notifications, notification)
	}
	return notifications, nil
}

func (u *UserRowType) UpdateAllNotificationsAsRead() error {
	stmt, err := db.Prepare(`UPDATE notifications SET "read" = 1 WHERE user_id = ?`)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	_, err = stmt.Exec((*u).Id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}
