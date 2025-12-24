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
		return removeLikeFromPost(userId, postId)
	}
	existingDislikeId, err := checkIfUserDislikedPost(userId, postId)
	if err != nil {
		return err
	}
	if existingDislikeId != 0 {
		return removeReaction(existingDislikeId)
	}
	return addLikeToPost(userId, postId)
}

func DoLikeComment(userId, commentId int64) error {
	alreadyLiked, err := hasUserLikedComment(userId, commentId)
	if err != nil {
		return err
	}
	if alreadyLiked {
		return removeLikeFromComment(userId, commentId)
	}
	existingDislikeId, err := checkIfUserDislikedComment(userId, commentId)
	if err != nil {
		return err
	}
	if existingDislikeId != 0 {
		return removeReaction(existingDislikeId)
	}
	return addLikeToComment(userId, commentId)
}
