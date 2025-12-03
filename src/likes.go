package forum

type Likes []Like

type Like struct {
	PostId    int
	UserId    int
	CommentId int
}

// Handle the business logic for liking a post
func DoLike(userId, postId int) error {
	// Check if user already liked this post
	alreadyLiked, err := hasUserAlreadyLikedPost(userId, postId)
	if err != nil {
		return err
	}

	if alreadyLiked {
		// User already liked the post, no action needed
		return nil
	}

	// Check if user previously disliked this post
	existingDislikeId, err := checkIfUserAlreadyDislikedPost(userId, postId)
	if err != nil {
		return err
	}

	if existingDislikeId != 0 {
		// User previously disliked, remove the dislike first
		err = removeDislikeFromPost(existingDislikeId)
		if err != nil {
			return err
		}
	}

	// Add the like
	err = addLikeToPost(userId, postId)
	if err != nil {
		return err
	}

	return nil
}

// Handle the business logic for unliking a post
func UndoLike(userId, postId int) error {
	return removeLikeFromPost(userId, postId)
}

// Handle the business logic for liking a comment
func DoLikeComment(userId, commentId int) error {
	// Check if user already liked this comment
	alreadyLiked, err := hasUserAlreadyLikedComment(userId, commentId)
	if err != nil {
		return err
	}

	if alreadyLiked {
		// User already liked the comment, no action needed
		return nil
	}

	// Check if user previously disliked this comment
	existingDislikeId, err := checkIfUserAlreadyDislikedComment(userId, commentId)
	if err != nil {
		return err
	}

	if existingDislikeId != 0 {
		// User previously disliked, remove the dislike first
		err = removeDislikeFromComment(existingDislikeId)
		if err != nil {
			return err
		}
	}

	// Add the like
	err = addLikeToComment(userId, commentId)
	if err != nil {
		return err
	}

	return nil
}

// Handle the business logic for unliking a comment
func UndoLikeComment(userId, commentId int) error {
	return removeLikeFromComment(userId, commentId)
}
