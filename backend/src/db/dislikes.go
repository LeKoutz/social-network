package db

import (
	"forum/src/ferror"
)

func SelectDislikesCountByPostId(postId int64) (int64, error) {
	var dislikes int64
	err := db.QueryRow(`
        SELECT COUNT(*)
        FROM reactions
        WHERE post_id = ? AND value = 2
    `, postId).Scan(&dislikes)
	if err != nil {
		return 0, ferror.ReturnErr(err)
	}
	return dislikes, nil
}

func SelectDislikesCountByCommentId(commentId int64) (int64, error) {
	var dislikes int64
	err := db.QueryRow(`
        SELECT COUNT(*)
        FROM reactions
        WHERE comment_id = ? AND value = 2
    `, commentId).Scan(&dislikes)
	if err != nil {
		return 0, ferror.ReturnErr(err)
	}
	return dislikes, nil
}
