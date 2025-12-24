package forum

type Dislikes []Dislike

type Dislike struct {
	PostId    int64
	UserId    int64
	CommentId int64
}

func DoDislikePost(userId, postId int64) error {
	alreadyDisliked, err := hasUserDislikedPost(userId, postId)
	if err != nil {
		return err
	}
	if alreadyDisliked {
		return UndoDislike(userId, postId)
	}
	existingLikeId, err := checkIfUserLikedPost(userId, postId)
	if err != nil {
		return err
	}
	if existingLikeId != 0 {
		err = removeLikeFromPost(userId, postId)
		if err != nil {
			return err
		}
	}
	err = addDislikeToPost(userId, postId)
	if err != nil {
		return err
	}
	return nil
}

func UndoDislike(userId, postId int64) error {
	dislikeId, err := checkIfUserDislikedPost(userId, postId)
	if err != nil {
		return err
	}
	if dislikeId == 0 {
		return nil
	}
	return removeDislikeFromPost(dislikeId)
}

func DoDislikeComment(userId, commentId int64) error {
	alreadyDisliked, err := hasUserDislikedComment(userId, commentId)
	if err != nil {
		return err
	}
	if alreadyDisliked {
		return UndoDislikeComment(userId, commentId)
	}
	existingLikeId, err := checkIfUserLikedComment(userId, commentId)
	if err != nil {
		return err
	}
	if existingLikeId != 0 {
		err = removeLikeFromComment(userId, commentId)
		if err != nil {
			return err
		}
	}
	err = addDislikeToComment(userId, commentId)
	if err != nil {
		return err
	}
	return nil
}

func UndoDislikeComment(userId, commentId int64) error {
	dislikeId, err := checkIfUserDislikedComment(userId, commentId)
	if err != nil {
		return err
	}
	if dislikeId == 0 {
		return nil
	}
	return removeDislikeFromComment(dislikeId)
}
