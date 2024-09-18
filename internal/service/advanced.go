package service

import (
	"forum/models"
	"forum/pkg/validator"
)

func (s *service) UpdatePost(data *models.TemplateData, form *models.PostForm) error {
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")

	if !form.Valid() {
		data.Form = form
		return models.ErrNotValidPostForm
	}
	err := s.repo.UpdateUserPost(data.PostID, form.Content, form.Title)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeletePost(postID int) error {
	return s.repo.DeleteUserPost(postID)
}
