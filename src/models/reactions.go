package models

import (
	"errors"
	"forum/src/db"
	"forum/src/utils"
)

type ReactionsType []ReactionType

func (reactions *ReactionsType) GetPostLikesByUserId(id int64) error {
	var err error
	var reaction_rows db.ReactionRowsType
	err = reaction_rows.SelectPostLikesByUserId(id)
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return err
	}
	for _, reaction_row := range reaction_rows {
		var reaction ReactionType
		reaction.FromReactionRowType(reaction_row)
		t, err := utils.ConvertStringToTime(reaction.TimestampString)
		if err != nil {
			if config.Debug {
				err = errors.Join(utils.GetFunctionName(), err)
			}
			return err
		}
		reaction.TimestampString = utils.ConvertTimeToString(t)
		*reactions = append(*reactions, reaction)
	}
	return nil
}

func GetPostDislikesByUserId(id int64) (ReactionsType, error) {
	var reactions ReactionsType
	reaction_rows, err := db.SelectPostDislikesByUserId(id)
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return ReactionsType{}, err
	}
	for _, reaction_row := range reaction_rows {
		var reaction ReactionType
		reaction.FromReactionRowType(reaction_row)
		t, err := utils.ConvertStringToTime(reaction.TimestampString)
		if err != nil {
			if config.Debug {
				err = errors.Join(utils.GetFunctionName(), err)
			}
			return ReactionsType{}, err
		}
		reaction.TimestampString = utils.ConvertTimeToString(t)
		reactions = append(reactions, reaction)
	}
	return reactions, nil
}

func GetCommentLikesByUserId(id int64) (ReactionsType, error) {
	var reactions ReactionsType
	reaction_rows, err := db.SelectCommentLikesByUserId(id)
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return ReactionsType{}, err
	}
	for _, reaction_row := range reaction_rows {
		var reaction ReactionType
		reaction.FromReactionRowType(reaction_row)
		t, err := utils.ConvertStringToTime(reaction.TimestampString)
		if err != nil {
			if config.Debug {
				err = errors.Join(utils.GetFunctionName(), err)
			}
			return ReactionsType{}, err
		}
		reaction.TimestampString = utils.ConvertTimeToString(t)
		reactions = append(reactions, reaction)
	}
	return reactions, nil
}

func GetCommentDisikesByUserId(id int64) (ReactionsType, error) {
	var reactions ReactionsType
	reaction_rows, err := db.SelectCommentDislikesByUserId(id)
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return ReactionsType{}, errors.Join(utils.GetFunctionName(), err)
	}
	for _, reaction_row := range reaction_rows {
		var reaction ReactionType
		reaction.FromReactionRowType(reaction_row)
		t, err := utils.ConvertStringToTime(reaction.TimestampString)
		if err != nil {
			if config.Debug {
				err = errors.Join(utils.GetFunctionName(), err)
			}
			return ReactionsType{}, err
		}
		reaction.TimestampString = utils.ConvertTimeToString(t)
		reactions = append(reactions, reaction)
	}
	return reactions, nil
}
