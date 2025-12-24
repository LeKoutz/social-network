package forum

type Likes []Like

type Like struct {
	PostId    int64
	UserId    int64
	CommentId int64
}

func DoLike(userId, postId int64) error {
	alreadyLiked, err := hasUserLikedPost(userId, postId)
	if err != nil {
		return err
	}
	if alreadyLiked {
		return UndoLike(userId, postId)
	}
	existingDislikeId, err := checkIfUserDislikedPost(userId, postId)
	if err != nil {
		return err
	}
	if existingDislikeId != 0 {
		err = removeDislikeFromPost(existingDislikeId)
		if err != nil {
			return err
		}
	}
	err = addLikeToPost(userId, postId)
	if err != nil {
		return err
	}
	return nil
}

func UndoLike(userId, postId int64) error {
	return removeLikeFromPost(userId, postId)
}

func DoLikeComment(userId, commentId int64) error {
	alreadyLiked, err := hasUserLikedComment(userId, commentId)
	if err != nil {
		return err
	}
	if alreadyLiked {
		return UndoLikeComment(userId, commentId)
	}
	existingDislikeId, err := checkIfUserDislikedComment(userId, commentId)
	if err != nil {
		return err
	}
	if existingDislikeId != 0 {
		err = removeDislikeFromComment(existingDislikeId)
		if err != nil {
			return err
		}
	}
	err = addLikeToComment(userId, commentId)
	if err != nil {
		return err
	}
	return nil
}

func UndoLikeComment(userId, commentId int64) error {
	return removeLikeFromComment(userId, commentId)
}
