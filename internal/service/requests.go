package service

import (
	"forum/models"
)

func (s *service) ModerationRequest(id int) error {
	role, err := s.repo.GetRole(id)
	if err != nil {
		return err
	}
	if role != "user" {
		return models.NoPermission
	}
	err = s.repo.CreateRequest(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) RequestsList() ([]*models.Request, error) {
	requests, err := s.repo.GetRequests()
	if err != nil {
		return nil, err
	}
	return requests, nil
}
