package controllers

import (
	"github.com/stretchr/testify/assert"
	"redditDataCompiler/models"
	"testing"
)

func Test_stats_CalculateTopUserAndPost(t *testing.T) {
	type fields struct {
		MostUpvotedPost   models.Post
		UserWithMostPosts models.User
		UserToTotalPosts  map[string]uint
		SeenPosts         map[string]struct{}
	}
	type args struct {
		response *models.APIResponse
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		expectedPost models.Post
		expectedUser models.User
	}{
		{
			name: "adding to empty stats",
			fields: fields{
				MostUpvotedPost:   models.Post{},
				UserWithMostPosts: models.User{},
				UserToTotalPosts:  map[string]uint{},
				SeenPosts:         map[string]struct{}{},
			},
			args: args{
				response: &models.APIResponse{
					Data: models.APIData{
						Children: []models.APIChildren{
							{
								Data: models.Listing{
									Name:           "sunnydale_high_123",
									Ups:            2,
									AuthorFullName: "buffy_summers_123",
									Author:         "buffy_summers",
									Title:          "sunnydale high",
								},
							},
							{
								Data: models.Listing{
									Name:           "the_bronze_123",
									Ups:            3,
									AuthorFullName: "willow_rosenberg_123",
									Author:         "willow_rosenberg",
									Title:          "the bronze",
								},
							},
							{
								Data: models.Listing{
									Name:           "uc_sunnydale_123",
									Ups:            1,
									AuthorFullName: "willow_rosenberg_123",
									Author:         "willow_rosenberg",
									Title:          "uc sunnydale",
								},
							},
						},
					},
				},
			},
			expectedPost: models.Post{
				Title:      "the bronze",
				NumUpvotes: 3,
				PostID:     "the_bronze_123",
			},
			expectedUser: models.User{
				AuthorID:       "willow_rosenberg_123",
				AuthorUserName: "willow_rosenberg",
				NumPosts:       2,
			},
		},
		{
			name: "replacing existing stat",
			fields: fields{
				MostUpvotedPost: models.Post{
					Title:      "the bronze",
					NumUpvotes: 3,
					PostID:     "the_bronze_123",
				},
				UserWithMostPosts: models.User{
					AuthorID:       "willow_rosenberg_123",
					AuthorUserName: "willow_rosenberg",
					NumPosts:       1,
				},
				UserToTotalPosts: map[string]uint{
					"willow_rosenberg_123": 1,
				},
				SeenPosts: map[string]struct{}{
					"the_bronze_123": {},
				},
			},
			args: args{
				response: &models.APIResponse{
					Data: models.APIData{
						Children: []models.APIChildren{
							{
								Data: models.Listing{
									Name:           "sunnydale_high_123",
									Ups:            5,
									AuthorFullName: "buffy_summers_123",
									Author:         "buffy_summers",
									Title:          "sunnydale high",
								},
							},
							{
								Data: models.Listing{
									Name:           "uc_sunnydale_123",
									Ups:            1,
									AuthorFullName: "buffy_summers_123",
									Author:         "buffy_summers",
									Title:          "uc sunnydale",
								},
							},
						},
					},
				},
			},
			expectedPost: models.Post{
				Title:      "sunnydale high",
				NumUpvotes: 5,
				PostID:     "sunnydale_high_123",
			},
			expectedUser: models.User{
				AuthorID:       "buffy_summers_123",
				AuthorUserName: "buffy_summers",
				NumPosts:       2,
			},
		},
		{
			name: "don't replacing existing stat",
			fields: fields{
				MostUpvotedPost: models.Post{
					Title:      "the bronze",
					NumUpvotes: 3,
					PostID:     "the_bronze_123",
				},
				UserWithMostPosts: models.User{
					AuthorID:       "willow_rosenberg_123",
					AuthorUserName: "willow_rosenberg",
					NumPosts:       2,
				},
				UserToTotalPosts: map[string]uint{
					"willow_rosenberg_123": 2,
				},
				SeenPosts: map[string]struct{}{
					"the_bronze_123": {},
				},
			},
			args: args{
				response: &models.APIResponse{
					Data: models.APIData{
						Children: []models.APIChildren{
							{
								Data: models.Listing{
									Name:           "sunnydale_high_123",
									Ups:            1,
									AuthorFullName: "buffy_summers_123",
									Author:         "buffy_summers",
									Title:          "sunnydale high",
								},
							},
						},
					},
				},
			},
			expectedPost: models.Post{
				Title:      "the bronze",
				NumUpvotes: 3,
				PostID:     "the_bronze_123",
			},
			expectedUser: models.User{
				AuthorID:       "willow_rosenberg_123",
				AuthorUserName: "willow_rosenberg",
				NumPosts:       2,
			},
		},
		{
			name: "update existing post upvote count",
			fields: fields{
				MostUpvotedPost: models.Post{
					Title:      "the bronze",
					NumUpvotes: 3,
					PostID:     "the_bronze_123",
				},
				UserWithMostPosts: models.User{
					AuthorID:       "willow_rosenberg_123",
					AuthorUserName: "willow_rosenberg",
					NumPosts:       2,
				},
				UserToTotalPosts: map[string]uint{
					"willow_rosenberg_123": 2,
				},
				SeenPosts: map[string]struct{}{
					"the_bronze_123": {},
				},
			},
			args: args{
				response: &models.APIResponse{
					Data: models.APIData{
						Children: []models.APIChildren{
							{
								Data: models.Listing{
									Name:           "the_bronze_123",
									Ups:            5,
									AuthorFullName: "willow_rosenberg_123",
									Author:         "willow_rosenberg",
									Title:          "the bronze",
								},
							},
						},
					},
				},
			},
			expectedPost: models.Post{
				Title:      "the bronze",
				NumUpvotes: 5,
				PostID:     "the_bronze_123",
			},
			expectedUser: models.User{
				AuthorID:       "willow_rosenberg_123",
				AuthorUserName: "willow_rosenberg",
				NumPosts:       2,
			},
		},
		{
			name: "update existing user post count",
			fields: fields{
				MostUpvotedPost: models.Post{
					Title:      "the bronze",
					NumUpvotes: 5,
					PostID:     "the_bronze_123",
				},
				UserWithMostPosts: models.User{
					AuthorID:       "willow_rosenberg_123",
					AuthorUserName: "willow_rosenberg",
					NumPosts:       2,
				},
				UserToTotalPosts: map[string]uint{
					"willow_rosenberg_123": 2,
				},
				SeenPosts: map[string]struct{}{
					"the_bronze_123": {},
				},
			},
			args: args{
				response: &models.APIResponse{
					Data: models.APIData{
						Children: []models.APIChildren{
							{
								Data: models.Listing{
									Name:           "the_library_123",
									Ups:            1,
									AuthorFullName: "willow_rosenberg_123",
									Author:         "willow_rosenberg",
									Title:          "the library",
								},
							},
						},
					},
				},
			},
			expectedPost: models.Post{
				Title:      "the bronze",
				NumUpvotes: 5,
				PostID:     "the_bronze_123",
			},
			expectedUser: models.User{
				AuthorID:       "willow_rosenberg_123",
				AuthorUserName: "willow_rosenberg",
				NumPosts:       3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &stats{
				MostUpvotedPost:   tt.fields.MostUpvotedPost,
				UserWithMostPosts: tt.fields.UserWithMostPosts,
				UserToTotalPosts:  tt.fields.UserToTotalPosts,
				SeenPosts:         tt.fields.SeenPosts,
			}
			s.CalculateTopUserAndPost(tt.args.response)
			assert.Equal(t, tt.expectedPost, s.MostUpvotedPost)
			assert.Equal(t, tt.expectedUser, s.UserWithMostPosts)
		})
	}
}
