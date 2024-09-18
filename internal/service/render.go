package service

import (
	"fmt"
	"forum/models"
	"strconv"
)

func (s service) MainLoader(data *models.TemplateData) (*models.TemplateData, error) {
	posts, err := s.LatestPosts()
	if err != nil {
		return nil, err
	}
	data.Posts = posts
	return data, nil
}

func (s service) HomeUpdates(data *models.TemplateData, filterData *models.FilterData, postID string, isStandardPath bool) (*models.TemplateData, error) {
	if filterData.Filter != "" {
		posts, err := s.LatestPosts()
		if err != nil {
			return nil, err
		}
		posts = filter(posts, filterData)
		data.Posts = posts
		return data, nil
	}

	intPostID, err := strconv.Atoi(postID)
	if err != nil && intPostID != 0 {
		return nil, err
	}

	if filterData.Reaction != "" {
		if filterData.Reaction == "1" {
			err := s.ReactionDone(data.UserID, intPostID, 1)
			if err != nil {
				return nil, err
			}
		} else if filterData.Reaction == "-1" {
			err := s.ReactionDone(data.UserID, intPostID, -1)
			if err != nil {
				return nil, err
			}
		}
	}

	posts, err := s.LatestPosts()
	if err != nil {
		return nil, err
	}
	if filterData.Filter != "" || !isStandardPath {
		posts = filter(posts, filterData)
	}

	data.Posts = posts

	return data, nil
}

func filter(posts []*models.Posts, filterData *models.FilterData) []*models.Posts {
	var filter string

	if filterData.Category != "" {
		filter = filterData.Category
	} else {
		filter = filterData.CategoryURL
	}

	var res []*models.Posts
	if filter != "All" {
		for _, v := range posts {
			for _, c := range v.Category {
				if c == filter {
					v.URL = filter
					res = append(res, v)
				}
			}
		}
	} else if filter == "All" {
		return posts
	}

	return res
}

func (s *service) LikedPostsLoader(data *models.TemplateData) (*models.TemplateData, error) {
	username, err := s.repo.GetUser(data.UserID)
	if err != nil {
		return nil, err
	}
	posts, err := s.repo.GetLikedPost(data.UserID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	data.Posts = posts
	data.Username = username

	return data, nil
}

func (s *service) DislikedPostsLoader(data *models.TemplateData) (*models.TemplateData, error) {
	username, err := s.repo.GetUser(data.UserID)
	if err != nil {
		return nil, err
	}
	posts, err := s.repo.GetDislikedPost(data.UserID)
	if err != nil {
		return nil, err
	}
	data.Posts = posts
	data.Username = username

	return data, nil
}

func (s *service) LikedPostsUpdate(data *models.TemplateData, postID int, reaction string) (*models.TemplateData, error) {
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

	data, err := s.LikedPostsLoader(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *service) DislikedPostsUpdate(data *models.TemplateData, postID int, reaction string) (*models.TemplateData, error) {
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

	data, err := s.LikedPostsLoader(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
