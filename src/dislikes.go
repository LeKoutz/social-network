package forum

type Dislikes []Dislike

type Dislike struct {
	PostId    int
	UserId    int
	CommentId int
}

// Handle the business logic for disliking a post
func DoDislikePost(userId, postId int) error {
	// Check if user already disliked this post
	alreadyDisliked, err := hasUserAlreadyDislikedPost(userId, postId)
	if err != nil {
		return err
	}

	if alreadyDisliked {
		// User already disliked the post, no action needed
		return nil
	}

	// Check if user previously liked this post
	existingLikeId, err := checkIfUserAlreadyLikedPost(userId, postId)
	if err != nil {
		return err
	}

	if existingLikeId != 0 {
		// User previously liked, remove the like first
		err = removeLikeFromPost(userId, postId)
		if err != nil {
			return err
		}
	}

	// Add the dislike
	err = addDislikeToPost(userId, postId)
	if err != nil {
		return err
	}

	return nil
}

// Handle the business logic for undisliking a post
func UndoDislike(userId, postId int) error {
	// First get the dislike ID
	dislikeId, err := checkIfUserAlreadyDislikedPost(userId, postId)
	if err != nil {
		return err
	}
	if dislikeId == 0 {
		return nil // No dislike to remove
	}
	return removeDislikeFromPost(dislikeId)
}

// Handle the business logic for disliking a comment
func DoDislikeComment(userId, commentId int) error {
	// Check if user already disliked this comment
	alreadyDisliked, err := hasUserAlreadyDislikedComment(userId, commentId)
	if err != nil {
		return err
	}

	if alreadyDisliked {
		// User already disliked the comment, no action needed
		return nil
	}

	// Check if user previously liked this comment
	existingLikeId, err := checkIfUserAlreadyLikedComment(userId, commentId)
	if err != nil {
		return err
	}

	if existingLikeId != 0 {
		// User previously liked, remove the like first
		err = removeLikeFromComment(userId, commentId)
		if err != nil {
			return err
		}
	}

	// Add the dislike
	err = addDislikeToComment(userId, commentId)
	if err != nil {
		return err
	}

	return nil
}

// Handle the business logic for undisliking a comment
func UndoDislikeComment(userId, commentId int) error {
	// First get the dislike ID
	dislikeId, err := checkIfUserAlreadyDislikedComment(userId, commentId)
	if err != nil {
		return err
	}
	if dislikeId == 0 {
		return nil // No dislike to remove
	}
	return removeDislikeFromComment(dislikeId)
}
