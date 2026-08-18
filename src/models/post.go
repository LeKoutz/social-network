package models

import (
	"errors"
	"forum/src/db"
	"forum/src/ferror"
	"forum/src/utils"
)

type PostType struct {
	db.PostRowType

	User       UserType
	Timestamp  int64
	Likes      int64
	Liked      bool
	Dislikes   int64
	Disliked   bool
	Category   CategoryType
	Categories CategoriesType
	Comments   CommentsType
}

func (p *PostType) ValidatePost() error {
	if len(p.Title) == 0 {
		return ferror.ErrorPostTitleEmpty
	}
	if len(p.Body) == 0 {
		return ferror.ErrorPostBodyEmpty
	}
	if p.Categories.IsEmpty() {
		return ferror.ErrorPostHasNoCategory
	}
	return nil
}

// Adds a Post in the database. Returns its id or error
func (p *PostType) Add() error {
	err := p.ValidatePost()
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return err
	}
	err = p.InsertPost()
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return err
	}
	for _, category := range p.Categories {
		err = db.AddPostCategory(db.PostCategoryRow{
			PostId:     p.Id,
			CategoryId: category.Id,
		})
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return errors.Join(utils.GetFunctionName(), err)
		}
	}
	return nil
}

func (p *PostType) GetReactions() error {
	var err error
	p.Likes, err = db.GetLikesCountByPostId(p.Id)
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	p.Dislikes, err = db.GetDislikesCountByPostId(p.Id)
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	return nil
}

func (p *PostType) GetReactionsByUserId(user_id int64) error {
	var err error
	p.Liked, err = HasUserLikedPost(user_id, p.Id)
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	p.Disliked, err = HasUserDislikedPost(user_id, p.Id)
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	return nil
}

func (p *PostType) GetComments() error {
	rows, err := p.SelectCommentsAndUsernameByPostId()
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return err
	}
	for _, row := range rows {
		var comment CommentType
		comment.CommentRowType = row
		t, err := utils.ConvertStringToTime(row.TimestampString)
		if err != nil {
			if config.Debug {
				err = errors.Join(utils.GetFunctionName(), err)
			}
			return err
		}
		comment.TimestampString = utils.ConvertTimeToString(t)
		p.Comments = append(p.Comments, comment)
	}
	return nil
}

func (p *PostType) Delete() error {
	err := p.DeletePostById()
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return err
	}
	return nil
}

func (p *PostType) Update() error {
	err := p.ValidatePost()
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return err
	}
	err = p.UpdatePost()
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return err
	}
	err = db.DeletePostCategoryByPostId(p.Id)
	if err != nil {
		if config.Debug {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		return err
	}
	for _, category := range p.Categories {
		err = db.AddPostCategory(db.PostCategoryRow{
			PostId:     p.Id,
			CategoryId: category.Id,
		})
		if err != nil {
			if config.Debug {
				err = errors.Join(utils.GetFunctionName(), err)
			}
			return err
		}
	}
	return nil
}

func (p *PostType) GetById() error {
	var err error
	err = p.SelectPostById()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	t, err := utils.ConvertStringToTime(p.TimestampString)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	p.TimestampString = utils.ConvertTimeToString(t)
	return nil
}
