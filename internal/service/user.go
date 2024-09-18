package service

import (
	"errors"
	"forum/models"
	"forum/pkg/validator"
)

func (s *service) UserSignUp(data *models.TemplateData, form *models.UserSignupForm) (*models.TemplateData, error) {
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")
	form.CheckField(validator.MaxChars(form.Password, 128), "password", "This field must be at most 128 characters long")
	form.CheckField(validator.MaxChars(form.Name, 128), "name", "This field must be at most 128 characters long")
	form.CheckField(validator.MaxChars(form.Email, 128), "email", "This field must be at most 128 characters long")
	form.CheckField(validator.Matches(form.Name, validator.StandardASCIIRX), "name", "Only standard ASCII characters are allowed")

	if !form.Valid() {
		data.Form = form
		return data, models.ErrNotValidPostForm
	}

	if form.Password != form.PasswordAgain {
		form.AddFieldError("password", "Password is incorrect")
		data.Form = form
		return data, models.ErrNotValidPostForm
	}
	id, err := s.repo.CreateUser(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			data.Form = form
			return data, models.ErrNotValidPostForm
		} else if errors.Is(err, models.ErrDuplicateName) {
			form.AddFieldError("name", "Usernmae is already in use")
			data.Form = form
			return data, models.ErrNotValidPostForm
		} else {
			return nil, err
		}
	}

	err = s.repo.CreateUserRole(id, "user")
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *service) UserLogin(data *models.TemplateData, form *models.UserLoginForm) (*models.TemplateData, int, error) {
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data.Form = form
		return data, 0, models.ErrNotValidPostForm
	}

	id, err := s.repo.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddFieldError("email", "Email or password is incorrect")

			data.Form = form

			return data, 0, models.ErrNotValidPostForm
		} else {
			return nil, 0, err
		}
	}
	return data, id, nil
}

func (s *service) SaveUser(form *models.UserSignupForm) (int, error) {
	// form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	// form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	// form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	// form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	// form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")
	// form.CheckField(validator.MaxChars(form.Password, 128), "password", "This field must be at most 128 characters long")
	// form.CheckField(validator.MaxChars(form.Name, 128), "name", "This field must be at most 128 characters long")
	// form.CheckField(validator.MaxChars(form.Email, 128), "email", "This field must be at most 128 characters long")
	// form.CheckField(validator.Matches(form.Name, validator.StandardASCIIRX), "name", "Only standard ASCII characters are allowed")

	// if !form.Valid() {
	// 	return 0, models.ErrNotValidPostForm
	// }

	id, err := s.repo.CreateUser(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			return 0, models.ErrDuplicateEmail
		} else if errors.Is(err, models.ErrDuplicateName) {
			form.AddFieldError("name", "Usernmae is already in use")
			return 0, models.ErrDuplicateName
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (s *service) GetUserByEmail(email string) (*models.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *service) GetUserData(data *models.TemplateData) (*models.TemplateData, error) {
	role, err := s.repo.GetRole(data.UserID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			return data, nil
		}
		return nil, err
	}
	data.Username, err = s.repo.GetUser(data.UserID)
	if err != nil {
		return nil, err
	}
	// potom chtoto budet

	data.Role = role
	return data, nil
}
