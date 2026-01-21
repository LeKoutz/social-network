package models

func RemoveReaction(reactionId int64) error {
	_, err := DB.Exec(`
		DELETE FROM reactions
		WHERE id = ?
		`, reactionId)
	return err
}
