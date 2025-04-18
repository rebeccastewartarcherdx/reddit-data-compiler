package controllers

import (
	"context"
	"errors"
	"fmt"
	"redditDataCompiler/models"
	"sync"
)

type Stats interface {
	CalculateTopUserAndPost(response *models.APIResponse)
	GetUserWithMostPosts(ctx context.Context) (*models.User, error)
	GetPostWithMostUpvotes(ctx context.Context) (*models.Post, error)
}
type stats struct {
	mu              sync.Mutex
	MostUpvotedPost models.Post

	UserWithMostPosts models.User
	UserToTotalPosts  map[string]uint
	SeenPosts         map[string]struct{}
}

func NewStatsProcessor() Stats {
	return &stats{
		UserToTotalPosts: map[string]uint{},
		SeenPosts:        map[string]struct{}{},
	}
}

func (s *stats) CalculateTopUserAndPost(response *models.APIResponse) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, r := range response.Data.Children {
		// update most upvoted post if needed
		if r.Data.Ups > s.MostUpvotedPost.NumUpvotes {
			s.MostUpvotedPost.NumUpvotes = r.Data.Ups
			s.MostUpvotedPost.PostID = r.Data.Name
			s.MostUpvotedPost.Title = r.Data.Title
		}

		// update user post counts/user with most posts if needed
		if _, ok := s.SeenPosts[r.Data.Name]; !ok {
			s.UserToTotalPosts[r.Data.AuthorFullName] += 1
			s.SeenPosts[r.Data.Name] = struct{}{}
		}
		if s.UserToTotalPosts[r.Data.AuthorFullName] > s.UserWithMostPosts.NumPosts {
			s.UserWithMostPosts.NumPosts = s.UserToTotalPosts[r.Data.AuthorFullName]
			s.UserWithMostPosts.AuthorUserName = r.Data.Author
			s.UserWithMostPosts.AuthorID = r.Data.AuthorFullName
		}
	}

	fmt.Println("most upvoted post: ", s.MostUpvotedPost)
	fmt.Println("user with most posts: ", s.UserWithMostPosts)

	return
}

func (s *stats) GetUserWithMostPosts(ctx context.Context) (*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.UserWithMostPosts.NumPosts == 0 {
		return nil, errors.New("no data compiled yet")
	}
	return &s.UserWithMostPosts, nil
}

func (s *stats) GetPostWithMostUpvotes(ctx context.Context) (*models.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.MostUpvotedPost.Title == "" {
		return nil, errors.New("no data compiled yet")
	}
	return &s.MostUpvotedPost, nil
}
