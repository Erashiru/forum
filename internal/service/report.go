package service

import "forum/models"

func (s *service) CreateReport(report *models.Report) error {
	return s.repo.CreateReport(report)
}

func (s *service) ReportsList() ([]*models.Report, error) {
	return s.repo.GetReports()
}

func (s *service) ReportDone(id int) error {
	return s.repo.ReportDone(id)
}
