package models

import (
	"errors"
	utils "forum/src/utils"
)

type Dislikes []Dislike

func getDislikesCountByCommentId(commentId int64) (int64, error) {
	var dislikes int64
	err := DB.QueryRow(`
        SELECT COUNT(*)
        FROM reactions
        WHERE comment_id = ? AND value = 2
    `, commentId).Scan(&dislikes)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return 0, err
	}
	return dislikes, nil
}

func getDislikesCountByPostId(postId int64) (int, error) {
	var dislikes int
	err := DB.QueryRow(`
        SELECT COUNT(*)
        FROM reactions
        WHERE post_id = ? AND value = 2
    `, postId).Scan(&dislikes)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return 0, err
	}
	return dislikes, nil
}
