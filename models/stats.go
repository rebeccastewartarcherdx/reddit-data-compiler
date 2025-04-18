package models

type Post struct {
	Title      string
	NumUpvotes uint
	PostID     string
}

type User struct {
	AuthorID       string
	AuthorUserName string
	NumPosts       uint
}
