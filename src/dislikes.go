package forum

type Dislikes []Dislike

type Dislike struct {
	PostId    int64
	UserId    int64
	CommentId int64
}

func DoDislikePost(userId, postId int64) error {
	dislikeId, err := checkIfUserDislikedPost(userId, postId)
	if err != nil {
		return err
	}
	if dislikeId != 0 {
		return removeDislikeFromPost(dislikeId)
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
	return addDislikeToPost(userId, postId)
}

func DoDislikeComment(userId, commentId int64) error {
	dislikeId, err := checkIfUserDislikedComment(userId, commentId)
	if err != nil {
		return err
	}
	if dislikeId != 0 {
		return removeReaction(dislikeId)
	}
	existingLikeId, err := checkIfUserLikedComment(userId, commentId)
	if err != nil {
		return err
	}
	if existingLikeId != 0 {
		err = removeReaction(existingLikeId)
		if err != nil {
			return err
		}
	}
	return addDislikeToComment(userId, commentId)
}
