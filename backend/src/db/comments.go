package db

import "forum/src/ferror"

type CommentRowsType []CommentRowType

func SelectCommentsByUserId(id int64) (CommentRowsType, error) {
	var comments CommentRowsType
	rows, err := db.Query(`
	SELECT id, post_id, body, COALESCE(image_path, ''), timestamp, user_id
	FROM comments
	WHERE user_id = ?`, id)
	if err != nil {
		return comments, ferror.ReturnErr(err)
	}
	for rows.Next() {
		var comment CommentRowType
		err = rows.Scan(&comment.Id, &comment.PostId, &comment.Body, &comment.ImagePath, &comment.Timestamp, &comment.UserId)
		if err != nil {
			return comments, ferror.ReturnErr(err)
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
