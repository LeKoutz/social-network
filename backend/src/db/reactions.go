package db

import "forum/src/ferror"

type ReactionRowsType []ReactionRowType

func (reaction_rows *ReactionRowsType) SelectPostLikesByUserId(id int64) error {
	var err error
	rows, err := db.Query(`
	SELECT id, post_id, user_id, timestamp
	FROM reactions
	WHERE user_id = ? AND value=1 AND post_id IS NOT NULL
	`, id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	defer rows.Close()
	for rows.Next() {
		var reaction ReactionRowType
		err = rows.Scan(&reaction.Id, &reaction.PostId, &reaction.UserId, &reaction.Timestamp)
		if err != nil {
			return ferror.ReturnErr(err)
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
		return ReactionRowsType{}, ferror.ReturnErr(err)
	}
	defer rows.Close()
	for rows.Next() {
		var reaction ReactionRowType
		err = rows.Scan(&reaction.Id, &reaction.PostId, &reaction.UserId, &reaction.Timestamp)
		if err != nil {
			return ReactionRowsType{}, ferror.ReturnErr(err)
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
		return ReactionRowsType{}, ferror.ReturnErr(err)
	}
	defer rows.Close()
	for rows.Next() {
		var reaction ReactionRowType
		err = rows.Scan(&reaction.Id, &reaction.CommentId, &reaction.UserId, &reaction.Timestamp)
		if err != nil {
			return ReactionRowsType{}, ferror.ReturnErr(err)
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
		return ReactionRowsType{}, ferror.ReturnErr(err)
	}
	defer rows.Close()
	for rows.Next() {
		var reaction ReactionRowType
		err = rows.Scan(&reaction.Id, &reaction.CommentId, &reaction.UserId, &reaction.Timestamp)
		if err != nil {
			return ReactionRowsType{}, ferror.ReturnErr(err)
		}
		reactions = append(reactions, reaction)
	}
	return reactions, nil
}
