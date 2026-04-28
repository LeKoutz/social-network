package models

import (
	"errors"
	"forum/src/utils"
)

type Reactions []Reaction

func GetPostLikesByUserId(id int64) (Reactions, error) {
	var reactions Reactions
	rows, err := DB.Query(`
	SELECT id, post_id, user_id, timestamp
	FROM reactions
	WHERE user_id = ? AND value=1
	`, id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return Reactions{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var reaction Reaction
		var ts string
		err = rows.Scan(&reaction.Id, &reaction.PostId, &reaction.UserId, &ts)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return Reactions{}, err
		}
		t, err := utils.ConvertStringToTime(ts)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return Reactions{}, err
		}
		reaction.TimestampString = utils.ConvertTimeToString(t)
		reactions = append(reactions, reaction)
	}
	return reactions, nil
}