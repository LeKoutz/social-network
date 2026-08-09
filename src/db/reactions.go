package db

import (
	"errors"
	"forum/src/utils"
)

type ReactionRowsType []ReactionRowType

func (reaction_rows *ReactionRowsType) SelectPostLikesByUserId(id int64) error {
	// var reactions ReactionRowsType
	var err error
	rows, err := db.Query(`
	SELECT id, post_id, user_id, timestamp
	FROM reactions
	WHERE user_id = ? AND value=1 AND post_id IS NOT NULL
	`, id)
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	defer rows.Close()
	for rows.Next() {
		var reaction ReactionRowType
		err = rows.Scan(&reaction.Id, &reaction.PostId, &reaction.UserId, &reaction.TimestampString)
		if err != nil {
			return errors.Join(utils.GetFunctionName(), err)
		}
		*reaction_rows = append(*reaction_rows, reaction)
	}
	return nil
}

func SelectPostDislikesByUserId(id int64) (ReactionRowsType, error) {
	var reactions ReactionRowsType
	rows, err := db.Query(`
	SELECT id, post_id, user_id, timestamp
	FROM reactions
	WHERE user_id = ? AND value=2 AND post_id IS NOT NULL
	`, id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return ReactionRowsType{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var reaction ReactionRowType
		err = rows.Scan(&reaction.Id, &reaction.PostId, &reaction.UserId, &reaction.TimestampString)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return ReactionRowsType{}, err
		}
		reactions = append(reactions, reaction)
	}
	return reactions, nil
}

func SelectCommentLikesByUserId(id int64) (ReactionRowsType, error) {
	var reactions ReactionRowsType
	rows, err := db.Query(`
	SELECT id, comment_id, user_id, timestamp
	FROM reactions
	WHERE user_id = ? AND value=1 AND comment_id IS NOT NULL
	`, id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return ReactionRowsType{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var reaction ReactionRowType
		err = rows.Scan(&reaction.Id, &reaction.CommentId, &reaction.UserId, &reaction.TimestampString)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return ReactionRowsType{}, err
		}
		reactions = append(reactions, reaction)
	}
	return reactions, nil
}

func SelectCommentDislikesByUserId(id int64) (ReactionRowsType, error) {
	var reactions ReactionRowsType
	rows, err := db.Query(`
	SELECT id, comment_id, user_id, timestamp
	FROM reactions
	WHERE user_id = ? AND value=2 AND comment_id IS NOT NULL
	`, id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return ReactionRowsType{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var reaction ReactionRowType
		err = rows.Scan(&reaction.Id, &reaction.CommentId, &reaction.UserId, &reaction.TimestampString)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return ReactionRowsType{}, err
		}
		reactions = append(reactions, reaction)
	}
	return reactions, nil
}
