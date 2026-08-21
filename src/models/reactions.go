package models

import (
	"errors"
	"forum/src/db"
	"forum/src/utils"
)

type ReactionsType []ReactionType

func (r *ReactionsType) GetPostLikesByUserId(id int64) error {
	var err error
	var reactions ReactionsType
	var rows db.ReactionRowsType
	err = rows.SelectPostLikesByUserId(id)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	for _, row := range rows {
		var reaction ReactionType
		reaction.ReactionRowType = row
		reactions = append(reactions, reaction)
	}
	*r = reactions
	return nil
}

func GetPostDislikesByUserId(id int64) (ReactionsType, error) {
	var reactions ReactionsType
	rows, err := db.SelectPostDislikesByUserId(id)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return ReactionsType{}, err
	}
	for _, row := range rows {
		var reaction ReactionType
		reaction.ReactionRowType = row
		reactions = append(reactions, reaction)
	}
	return reactions, nil
}

func GetCommentLikesByUserId(id int64) (ReactionsType, error) {
	var reactions ReactionsType
	rows, err := db.SelectCommentLikesByUserId(id)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return ReactionsType{}, err
	}
	for _, row := range rows {
		var reaction ReactionType
		reaction.ReactionRowType = row
		reactions = append(reactions, reaction)
	}
	return reactions, nil
}

func GetCommentDisikesByUserId(id int64) (ReactionsType, error) {
	var reactions ReactionsType
	rows, err := db.SelectCommentDislikesByUserId(id)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return ReactionsType{}, err
	}
	for _, row := range rows {
		var reaction ReactionType
		reaction.ReactionRowType = row
		reactions = append(reactions, reaction)
	}
	return reactions, nil
}
