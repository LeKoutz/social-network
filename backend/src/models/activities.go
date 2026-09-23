package models

import (
	"forum/src/db"
	"forum/src/ferror"
	"sort"
)

type ActivitiesType []ActivityType

func (a *ActivitiesType) GetActivityById(id int64) error {
	err := a.GetPostsActivityById(id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	err = a.GetCommentsActivityById(id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	err = a.GetLikedPostsActivityById(id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	err = a.GetDislikedPostsActivityById(id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	err = a.GetLikedCommentsActivityById(id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	err = a.GetDislikedCommentsActivityById(id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	sort.Slice(*a, func(i, j int) bool {
		return (*a)[i].Timestamp > (*a)[j].Timestamp
	})
	return nil
}

func (a *ActivitiesType) GetPostsActivityById(id int64) error {
	var owner UserType
	owner.Id = id
	posts, err := owner.GetPosts()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	for _, post := range posts {
		err := post.SelectPostById()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		err = post.GetReactions()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		err = post.GetReactionsByUserId(id)
		if err != nil {
			return ferror.ReturnErr(err)
		}
		var activity ActivityType
		activity.Timestamp = post.Timestamp
		activity.Post = post
		activity.Type = "post"
		*a = append(*a, activity)
	}
	return nil
}

func (a *ActivitiesType) GetCommentsActivityById(id int64) error {
	rows, err := db.SelectCommentsByUserId(id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	for _, row := range rows {
		var activity ActivityType
		var post PostType
		var comment CommentType
		comment.CommentRowType = row
		post.Id = comment.PostId

		err = post.SelectPostById()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		err = comment.GetReactions()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		err = comment.GetReactionsByUserId(id)
		if err != nil {
			return ferror.ReturnErr(err)
		}
		activity.Type = "comment"
		activity.Comment = comment
		activity.Timestamp = comment.Timestamp
		activity.Post = post
		*a = append(*a, activity)
	}
	return nil
}

func (a *ActivitiesType) GetLikedPostsActivityById(id int64) error {
	var reactions ReactionsType
	err := reactions.GetPostLikesByUserId(id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	for _, reaction := range reactions {
		var activity ActivityType
		var post PostType
		post.Id = reaction.PostId
		err = post.SelectPostById()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		err = post.GetReactions()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		err = post.GetReactionsByUserId(id)
		if err != nil {
			return ferror.ReturnErr(err)
		}
		activity.Type = "postLike"
		activity.Timestamp = reaction.Timestamp
		activity.Post = post
		*a = append(*a, activity)
	}
	return nil
}

func (a *ActivitiesType) GetDislikedPostsActivityById(id int64) error {
	reactions, err := GetPostDislikesByUserId(id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	for _, reaction := range reactions {
		var activity ActivityType
		var post PostType
		post.Id = reaction.PostId
		err = post.SelectPostById()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		err = post.GetReactions()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		err = post.GetReactionsByUserId(id)
		if err != nil {
			return ferror.ReturnErr(err)
		}
		activity.Type = "postDislike"
		activity.Timestamp = reaction.Timestamp
		activity.Post = post
		*a = append(*a, activity)
	}
	return nil
}

func (a *ActivitiesType) GetLikedCommentsActivityById(id int64) error {
	reactions, err := GetCommentLikesByUserId(id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	for _, reaction := range reactions {
		var activity ActivityType
		var comment CommentType
		comment.Id = reaction.CommentId
		err = comment.SelectCommentById()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		err = comment.GetReactions()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		err = comment.GetReactionsByUserId(id)
		if err != nil {
			return ferror.ReturnErr(err)
		}
		var post PostType
		post.Id = comment.PostId
		err = post.SelectPostById()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		activity.Type = "commentLike"
		activity.Timestamp = reaction.Timestamp
		activity.Comment = comment
		activity.Post = post
		*a = append(*a, activity)
	}
	return nil
}

func (a *ActivitiesType) GetDislikedCommentsActivityById(id int64) error {
	reactions, err := GetCommentDisikesByUserId(id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	for _, reaction := range reactions {
		var activity ActivityType
		var comment CommentType
		comment.Id = reaction.CommentId
		err = comment.SelectCommentById()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		err = comment.GetReactions()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		err = comment.GetReactionsByUserId(id)
		if err != nil {
			return ferror.ReturnErr(err)
		}
		var post PostType
		post.Id = comment.PostId
		err = post.SelectPostById()
		if err != nil {
			return ferror.ReturnErr(err)
		}
		activity.Type = "commentDislike"
		activity.Timestamp = reaction.Timestamp
		activity.Comment = comment
		activity.Post = post
		*a = append(*a, activity)
	}
	return nil
}
