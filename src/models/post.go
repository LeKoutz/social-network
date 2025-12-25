package models

import (
	"errors"
	"forum/src/utils"
)

type Post struct {
	Id              int64
	Title           string
	Body            string
	UserId          int64
	User            User
	Timestamp       int64
	TimestampString string
	Likes           int
	Liked           bool
	Dislikes        int
	Disliked        bool
	Category        Category
	Categories      Categories
	Comments        Comments
}

func (p *Post) ValidatePost() error {
	if len(p.Title) == 0 {
		return ErrorPostTitleEmpty
	}
	if len(p.Body) == 0 {
		return ErrorPostBodyEmpty
	}
	if p.Categories.IsEmpty() {
		return ErrorPostHasNoCategory
	}
	return nil
}

func AddPost(post Post) (int64, error) {
	err := post.ValidatePost()
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return 0, err
	}
	stmt, err := DB.Prepare("INSERT INTO posts (title, body, user_id, timestamp) VALUES (?, ?, ?, ?)")
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return 0, err
	}
	res, err := stmt.Exec(post.Title, post.Body, post.UserId, utils.GetCurrentTimestamp())
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return 0, err
	}
	postId, err := res.LastInsertId()
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return 0, err
	}
	for _, category := range post.Categories {
		err = AddCategoryToPost(postId, category.Id)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return 0, err
		}
	}
	return postId, nil
}

func GetPostById(id int64) (Post, error) {
	var post Post
	var ts string
	err := DB.QueryRow(`SELECT title, body, timestamp, user_id FROM posts WHERE id = ?`, id).Scan(&post.Title, &post.Body, &ts, &post.UserId)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return Post{}, err
	}
	t, err := utils.ConvertStringToTime(ts)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return Post{}, err
	}
	post.TimestampString = utils.ConvertTimeToString(t)
	post.Id = id
	post.User, err = getUserById(post.UserId)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return Post{}, err
	}
	return post, nil
}

func (p *Post) GetReactions() error {
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

func (p *Post) GetReactionsByUserId(user_id int64) error {
	var err error
	(*p).Liked, err = HasUserLikedPost(user_id, (*p).Id)
	if err != nil {
		return err
	}
	(*p).Disliked, err = HasUserDislikedPost(user_id, (*p).Id)
	if err != nil {
		return err
	}
	return nil
}
