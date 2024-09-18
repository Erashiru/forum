package service

import (
	"forum/models"
)

func (s *service) UserPosts(data *models.TemplateData) (*models.TemplateData, error) {
	username, err := s.repo.GetUser(data.UserID)
	if err != nil {
		return nil, err
	}
	posts, err := s.repo.UserPosts(data.UserID)
	if err != nil {
		return nil, err
	}

	data.Posts = posts
	data.Username = username

	return data, nil
}

func (s *service) UserPostsUpdate(data *models.TemplateData, postID int, reaction string) (*models.TemplateData, error) {
	if reaction == "1" {
		err := s.ReactionDone(data.UserID, postID, 1)
		if err != nil {
			return nil, err
		}
	} else if reaction == "-1" {
		err := s.ReactionDone(data.UserID, postID, -1)
		if err != nil {
			return nil, err
		}
	}
	// }
	data, err := s.UserPosts(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
