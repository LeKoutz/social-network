package db

import "forum/src/ferror"

func SelectLikesCountByPostId(postId int64) (int64, error) {
	var likes int64
	err := db.QueryRow(`
        SELECT COUNT(*)
        FROM reactions
        WHERE post_id = ? AND value = 1
    `, postId).Scan(&likes)
	if err != nil {
		return 0, ferror.ReturnErr(err)
	}
	return likes, nil
}

func SelectLikesCountByCommentId(commentId int64) (int64, error) {
	var likes int64
	err := db.QueryRow(`
        SELECT COUNT(*)
        FROM reactions
        WHERE comment_id = ? AND value = 1
    `, commentId).Scan(&likes)
	if err != nil {
		return 0, ferror.ReturnErr(err)
	}
	return likes, nil
}
