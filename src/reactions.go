package forum

func (p *Post) getReactions() error {
	var err error
	(*p).Likes, err = getLikesCountByPostId((*p).Id)
	if err != nil {
		return err
	}
	(*p).Dislikes, err = getDislikesCountByPostId((*p).Id)
	if err != nil {
		return err
	}
	return nil
}

func (c *Comment) getReactions() error {
	var err error
	(*c).Likes, err = getLikesCountByCommentId((*c).Id)
	if err != nil {
		return err
	}
	(*c).Dislikes, err = getDislikesCountByCommentId((*c).Id)
	if err != nil {
		return err
	}
	return nil
}

func (p *Post) getReactionsByUserId(user_id int64) error {
	var err error
	(*p).Liked, err = hasUserLikedPost(user_id, (*p).Id)
	if err != nil {
		return err
	}
	(*p).Disliked, err = hasUserDislikedPost(user_id, (*p).Id)
	if err != nil {
		return err
	}
	return nil
}

func (c *Comment) getReactionsByUserId(user_id int64) error {
	var err error
	(*c).Liked, err = hasUserLikedComment(user_id, (*c).Id)
	if err != nil {
		return err
	}
	(*c).Disliked, err = hasUserDislikedComment(user_id, (*c).Id)
	if err != nil {
		return err
	}
	return nil
}
