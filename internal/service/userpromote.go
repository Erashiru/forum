package service

func (s *service) PromoteUser(id int, role string) error {
	return s.repo.UpdateUserRole(id, role)
}
