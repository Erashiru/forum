package service

import "forum/models"

func (s *service) UpdateUserRole(data *models.TemplateData, role string) error {
	if data.Role != "admin" {
		return models.NoPermission
	}

	return s.repo.UpdateUserRole(data.ReqUsrID, role)
}

func (s *service) PostDelete(postid int) error {
	err := s.repo.ReportDone(postid)
	if err != nil {
		return err
	}
	return s.repo.PostDelete(postid)
}
