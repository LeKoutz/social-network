package db

import (
	"errors"
	utils "forum/src/utils"
)

func GetLikesCountByPostId(postId int64) (int64, error) {
	var likes int64
	err := db.QueryRow(`
        SELECT COUNT(*)
        FROM reactions
        WHERE post_id = ? AND value = 1
    `, postId).Scan(&likes)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return 0, err
	}
	return likes, nil
}

func GetLikesCountByCommentId(commentId int64) (int64, error) {
	var likes int64
	err := db.QueryRow(`
        SELECT COUNT(*)
        FROM reactions
        WHERE comment_id = ? AND value = 1
    `, commentId).Scan(&likes)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return 0, err
	}
	return likes, nil
}
