package service

import (
	"database/sql"
	"errors"
	"forum/models"
	"forum/pkg/validator"
)

func (s service) LatestPosts() ([]*models.Posts, error) {
	posts, err := s.repo.LatestPosts()
	if err != nil {
		return nil, err
	}
	return posts, err
}

func (s service) ReactionDone(used_id, post_id, reaction int) error {
	return s.repo.CreateReaction(used_id, post_id, reaction)
}

func (s *service) GetPostData(data *models.TemplateData, postId int) (*models.TemplateData, error) {
	post, err := s.repo.GetPost(postId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	comments, err := s.repo.GetComments(postId)
	if err != nil {
		return nil, err
	}

	likes, err := s.repo.GetLikes(postId)
	if err != nil {
		return nil, err
	}
	dislikes, err := s.repo.GetDislikes(postId)
	if err != nil {
		return nil, err
	}
	category, err := s.repo.GetCategory(postId)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	data.Post = *post
	data.Post.Likes = likes
	data.Post.Dislikes = dislikes
	data.Post.Comments = comments
	data.Post.Category = category
	return data, nil
}

func (s *service) PostUpdate(data *models.TemplateData, userReactions *models.UserReaction, postID int) (*models.TemplateData, error) {
	data, err := s.GetPostData(data, postID)
	if err != nil {
		return nil, err
	}

	post, err := s.repo.GetPost(postID)
	if err != nil {
		if err == models.ErrNoRecord {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	data.Comments, err = s.repo.GetComments(postID)
	if err != nil {
		return nil, err
	}

	if userReactions.Reaction != "" {
		if userReactions.Reaction == "1" {
			err := s.repo.CreateReaction(data.UserID, postID, 1)
			if err != nil {
				return nil, err
			}
		} else if userReactions.Reaction == "-1" {
			err := s.repo.CreateReaction(data.UserID, postID, -1)
			if err != nil {
				return nil, err
			}
		}
	}

	if userReactions.CommentReaction != "" {
		if userReactions.CommentReaction == "1" {
			err := s.repo.CreateCommentReaction(data.UserID, userReactions.CommentID, 1)
			if err != nil {
				return nil, err
			}
		} else if userReactions.CommentReaction == "-1" {
			err := s.repo.CreateCommentReaction(data.UserID, userReactions.CommentID, -1)
			if err != nil {
				return nil, err
			}
		}
	}

	Form := &models.CommentForm{
		Text:   userReactions.Comment,
		Userid: data.UserID,
	}

	if Form.Text != "" {
		Form.CheckField(validator.NotBlank(Form.Text), "comment", "This field cannot be blank")
		Form.CheckField(validator.MaxChars(Form.Text, 100), "comment", "This field cannot be more than 100 characters long")

		if !Form.Valid() {
			data.Form = Form
			data.Post = *post
			return data, models.ErrNotValidPostForm
		}

		s.repo.CreateComment(postID, Form.Userid, Form.Text)
	}

	data.Form = Form
	data.Post = *post

	return data, nil
}

func (s *service) PostCreate(data *models.TemplateData, Form *models.PostForm) (*models.TemplateData, int, error) {
	permittedCategories := []string{"Python", "Golang", "JavaScript", "AI", "Algorithms"}

	Form.CheckField(validator.NotBlank(Form.Title), "title", "This field cannot be blank")
	Form.CheckField(validator.MaxChars(Form.Title, 100), "title", "This field cannot be more than 100 characters long")
	Form.CheckField(validator.NotBlank(Form.Content), "content", "This field cannot be blank")
	Form.CheckField(validator.PermittedCategories(Form.Category, permittedCategories), "category", "This category is not permitted")
	Form.CheckField(validator.NotBlankCato(Form.Category), "category", "This field cannot be blank")
	Form.CheckField(validator.MaxChars(Form.Content, 250), "content", "This field cannot be more than 250 characters long")
	Form.CheckField(validator.MaxChars(Form.Title, 250), "title", "This field cannot be more than 250 characters long")

	if !Form.Valid() {
		data.Form = Form
		return data, 0, models.ErrNotValidPostForm
	}

	id, err := s.repo.CreatePost(Form.Title, Form.Content, data.UserID)
	if err != nil {
		return nil, 0, err
	}

	err = s.repo.ChooseCategories(id, Form.Category)
	data.Form = Form
	return data, id, nil
}
