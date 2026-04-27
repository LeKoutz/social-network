package models

import (
	"errors"
	"forum/src/utils"
	"sort"
)

type Activity struct {
	Type	  string
	Timestamp string
	Post	  Post
	Comment   Comment
}

func (u *User) GetActivity() error {
	err := u.GetPostsActivity()
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	err = u.GetCommentsActivity()
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	err = u.GetLikedPostsActivity()
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	err = u.GetDislikedPostsActivity()
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	err = u.GetLikedCommentsActivity()
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	err = u.GetDislikedCommentsActivity()
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	sort.Slice(u.Activities, func(i, j int) bool {
		return u.Activities[i].Timestamp > u.Activities[j].Timestamp
	})
	return nil
}

func (u *User) GetPostsActivity() error {
	posts, err := u.GetPosts()
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
	}
	for _, post := range posts {
		err := post.GetById()
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		err = post.GetReactions()
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		err = post.GetReactionsByUserId(u.Id)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		var activity Activity
		activity.Timestamp = post.TimestampString
		activity.Post = post
		activity.Type = "post"
		u.Activities = append(u.Activities, activity)
	}
	return nil
}

func (u *User) GetCommentsActivity() error {
	rows, err := DB.Query(`
	SELECT id, post_id, body, timestamp, user_id
	FROM comments
	WHERE user_id = ?`, (*u).Id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	for rows.Next() {
		var activity Activity
		var comment Comment
		var ts string
		err = rows.Scan(&comment.Id, &comment.PostId, &comment.Body, &ts, &comment.UserId)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		t, err := utils.ConvertStringToTime(ts)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		rows.Close()
		err = comment.GetReactions()
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		err = comment.GetReactionsByUserId((*u).Id)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		var post = Post{Id: comment.PostId}
		err = post.GetById()
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		err = post.GetReactions()
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		err = post.GetReactionsByUserId((*u).Id)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		comment.TimestampString = utils.ConvertTimeToString(t)
		activity.Type = "comment"
		activity.Timestamp = comment.TimestampString
		activity.Comment = comment
		activity.Post = post
		u.Activities = append(u.Activities, activity)
	}
	return nil
}

func (u *User) GetLikedPostsActivity() error {
	rows, err := DB.Query(`
	SELECT p.id, p.title, p.body, p.user_id, r.timestamp
	FROM posts p
	JOIN reactions r ON p.id = r.post_id
	WHERE r.user_id = ? AND r.value = 1
	`, (*u).Id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	for rows.Next() {
		var activity Activity
		var post Post
		var ts string
		err = rows.Scan(&post.Id, &post.Title, &post.Body, &post.UserId, &ts)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		t, err := utils.ConvertStringToTime(ts)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		rows.Close()
		err = post.GetReactions()
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		err = post.GetReactionsByUserId((*u).Id)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		activity.Type = "postLike"
		activity.Timestamp = utils.ConvertTimeToString(t)
		activity.Post = post
		u.Activities = append(u.Activities, activity)
	}
	return nil
}

func (u *User) GetDislikedPostsActivity() error {
	rows, err := DB.Query(`
	SELECT p.id, p.title, p.body, p.user_id, r.timestamp
	FROM posts p
	JOIN reactions r ON p.id = r.post_id
	WHERE r.user_id = ? AND r.value = 2
	`, (*u).Id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	for rows.Next() {
		var activity Activity
		var post Post
		var ts string
		err = rows.Scan(&post.Id, &post.Title, &post.Body, &post.UserId, &ts)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		t, err := utils.ConvertStringToTime(ts)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		rows.Close()
		err = post.GetReactions()
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		err = post.GetReactionsByUserId((*u).Id)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		activity.Type = "postDislike"
		activity.Timestamp = utils.ConvertTimeToString(t)
		activity.Post = post
		u.Activities = append(u.Activities, activity)
	}
	return nil
}

func (u *User) GetLikedCommentsActivity() error {
	rows, err := DB.Query(`
	SELECT c.id, c.post_id, c.body, c.user_id, r.timestamp
	FROM comments c
	JOIN reactions r ON c.id = r.comment_id
	WHERE r.user_id = ? AND r.value = 1
	`, (*u).Id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	for rows.Next() {
		var activity Activity
		var comment Comment
		var ts string
		err = rows.Scan(&comment.Id, &comment.PostId, &comment.Body, &comment.UserId, &ts)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		t, err := utils.ConvertStringToTime(ts)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		rows.Close()
		err = comment.GetReactions()
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		err = comment.GetReactionsByUserId((*u).Id)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		var post = Post{Id: comment.PostId}
		err = post.GetById()
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		activity.Type = "commentLike"
		activity.Timestamp = utils.ConvertTimeToString(t)
		activity.Comment = comment
		activity.Post = post
		u.Activities = append(u.Activities, activity)
	}
	return nil
}

func (u *User) GetDislikedCommentsActivity() error {
	rows, err := DB.Query(`
	SELECT c.id, c.post_id, c.body, c.user_id, r.timestamp
	FROM comments c
	JOIN reactions r ON c.id = r.comment_id
	WHERE r.user_id = ? AND r.value = 2
	`, (*u).Id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	for rows.Next() {
		var activity Activity
		var comment Comment
		var ts string
		err = rows.Scan(&comment.Id, &comment.PostId, &comment.Body, &comment.UserId, &ts)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		t, err := utils.ConvertStringToTime(ts)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		rows.Close()
		err = comment.GetReactions()
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		err = comment.GetReactionsByUserId((*u).Id)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		var post = Post{Id: comment.PostId}
		err = post.GetById()
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		activity.Type = "commentDislike"
		activity.Timestamp = utils.ConvertTimeToString(t)
		activity.Comment = comment
		activity.Post = post
		u.Activities = append(u.Activities, activity)
	}
	return nil
}