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
		err := ferror.ErrorPostTitleEmpty
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	if len(p.Body) == 0 {
		err := ferror.ErrorPostBodyEmpty
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	if p.Categories.IsEmpty() {
		err := ferror.ErrorPostHasNoCategory
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}

// Adds a Post in the database. Returns its id or error
func (p *PostType) Add() error {
	err := p.ValidatePost()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	err = p.InsertPost()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	for _, category := range p.Categories {
		err = db.InsertPostCategory(db.PostCategoryRow{
			PostId:     p.Id,
			CategoryId: category.Id,
		})
		if err != nil {
			if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return err
		}
	}
	return nil
}

func (p *PostType) GetReactions() error {
	var err error
	p.Likes, err = db.SelectLikesCountByPostId(p.Id)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	p.Dislikes, err = db.SelectDislikesCountByPostId(p.Id)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}

func (p *PostType) GetReactionsByUserId(user_id int64) error {
	var err error
	p.Liked, err = HasUserLikedPost(user_id, p.Id)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	p.Disliked, err = HasUserDislikedPost(user_id, p.Id)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}

func (p *PostType) GetComments() error {
	rows, err := p.SelectCommentsAndUsernameByPostId()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	for _, row := range rows {
		var comment CommentType
		comment.CommentRowType = row
		p.Comments = append(p.Comments, comment)
	}
	return nil
}

func (p *PostType) Delete() error {
	err := p.DeletePostById()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}

func (p *PostType) Update() error {
	err := p.ValidatePost()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	err = p.UpdatePost()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	err = db.DeletePostCategoryByPostId(p.Id)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	for _, category := range p.Categories {
		err = db.InsertPostCategory(db.PostCategoryRow{
			PostId:     p.Id,
			CategoryId: category.Id,
		})
		if err != nil {
			if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
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
	return nil
}
